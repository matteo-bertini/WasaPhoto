# WasaPhoto

WasaPhoto è un social network per la condivisione di foto, sviluppato come progetto per il corso di **Web and Software Architectures** e successivamente ampliato come **tesi di laurea triennale in Ingegneria Informatica** presso Sapienza Università di Roma.

Il progetto implementa un'architettura client-server completa: un backend REST in **Go**, un frontend **Vue.js 3**, e uno storage persistente su **SQLite**. Le specifiche dell'API sono definite tramite **OpenAPI 3.0** (`doc/api.yaml`).

## Indice

- [Funzionalità](#funzionalità)
- [Stack tecnologico](#stack-tecnologico)
- [Struttura del progetto](#struttura-del-progetto)
- [Architettura](#architettura)
- [Modello dati](#modello-dati)
- [API REST](#api-rest)
- [Avvio del progetto](#avvio-del-progetto)
  - [Backend](#backend)
  - [Frontend](#frontend)
  - [Docker](#docker)
- [Configurazione](#configurazione)
- [Licenza](#licenza)

## Funzionalità

- **Autenticazione**: login/registrazione tramite un unico endpoint (`isSignUp`), sessione gestita con token Bearer.
- **Gestione profilo**: visualizzazione profilo con statistiche (follower, following, numero post), modifica username, eliminazione account.
- **Post**: upload foto con didascalia, cancellazione, recupero dell'immagine binaria.
- **Interazioni social**: like/unlike, commenti (creazione, lettura, cancellazione).
- **Relazioni sociali**: follow/unfollow, liste di follower e following.
- **Ban**: possibilità di bloccare altri utenti; un ban impedisce la visibilità reciproca di profilo, post, stream e liste social.
- **Stream**: feed paginato che privilegia i post degli utenti seguiti e completa gli slot rimanenti con post suggeriti, escludendo sempre gli utenti bannati.

## Stack tecnologico

**Backend**
- [Go](https://go.dev/) 1.25
- [httprouter](https://github.com/julienschmidt/httprouter) — routing HTTP
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — driver SQLite
- [logrus](https://github.com/sirupsen/logrus) — logging strutturato
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — hashing delle password
- [ksuid](https://github.com/segmentio/ksuid) — identificatori univoci ordinabili per utenti/post/sessioni
- [ardanlabs/conf](https://github.com/ardanlabs/conf) — configurazione da flag/env/YAML
- [gorilla/handlers](https://github.com/gorilla/handlers) — middleware CORS

**Frontend**
- [Vue.js 3](https://vuejs.org/)
- [Vue Router 4](https://router.vuejs.org/)
- [Axios](https://axios-http.com/)
- [Vite](https://vitejs.dev/) — build tool e dev server

**Storage**
- SQLite, con schema auto-inizializzato all'avvio e vincoli `FOREIGN KEY ... ON DELETE CASCADE` per la pulizia automatica dei dati correlati alla cancellazione di un account.

## Struttura del progetto

```
WasaPhoto/
├── cmd/
│   ├── webapi/            # entry point del server (main, config, CORS, web UI embedding)
│   └── healthcheck/       # binario di utilità per l'health check del container Docker
├── service/
│   ├── api/               # handler HTTP, routing, middleware di autenticazione
│   ├── database/          # layer di accesso ai dati (schema + query)
│   └── models/            # strutture dati condivise ed errori di dominio
├── webui/                 # applicazione frontend Vue.js
├── doc/
│   └── api.yaml           # specifica OpenAPI 3.0 delle API REST
├── demo/
│   └── config.yml          # esempio di file di configurazione
├── Dockerfile.backend      # build multi-stage del backend
└── open-npm.sh             # helper per lavorare sul frontend dentro un container Node
```

## Architettura

Il backend segue un'architettura a livelli con separazione netta delle responsabilità:

- **`service/api`**: espone gli endpoint REST, valida gli input, traduce gli errori di dominio in codici di stato HTTP. L'autenticazione è centralizzata in un middleware (`AuthMiddleware`) che valida il token Bearer e inietta l'`UserID` nel contesto della richiesta.
- **`service/database`**: incapsula tutta la logica SQL dietro l'interfaccia `AppDatabase`, così l'implementazione concreta (SQLite) resta disaccoppiata dagli handler.
- **`service/models`**: definisce le strutture condivise (Post, Comment, UserProfile, ecc.) e gli errori di dominio (es. `ErrUserNotFound`, `ErrForbiddenAction`) usati per mappare le risposte HTTP.

Le operazioni multi-step che devono restare atomiche (registrazione utente, like con controllo ban, follow, aggiunta commento) sono eseguite all'interno di transazioni SQL.

Il frontend è una Single Page Application Vue.js che comunica con il backend via REST, mantenendo la sessione (token e username) in `localStorage`.

## Modello dati

Lo schema SQLite viene creato automaticamente all'avvio (`service/database/database.go`) con le seguenti tabelle principali:

| Tabella     | Descrizione                                                       |
|-------------|---------------------------------------------------------------------|
| `accounts`  | Credenziali e token di sessione                                    |
| `profiles`  | Bio e avatar dell'utente                                           |
| `posts`     | Post pubblicati (didascalia, autore, data)                         |
| `follows`   | Relazioni di follow (follower → followed)                          |
| `likes`     | Like sui post                                                      |
| `comments`  | Commenti sui post                                                  |
| `bans`      | Relazioni di ban (banner → banned)                                 |

Tutte le tabelle collegate a `accounts` usano `ON DELETE CASCADE`, così la cancellazione di un account rimuove automaticamente profilo, post, like, commenti, follow e ban associati.

## API REST

Le API sono documentate integralmente in `doc/api.yaml` (OpenAPI 3.0). Endpoint principali:

| Metodo | Path | Descrizione |
|--------|------|-------------|
| `POST` | `/session` | Login o registrazione (`isSignUp`) |
| `DELETE` | `/session` | Logout |
| `GET` | `/users/{username}` | Profilo utente, statistiche e post |
| `PUT` | `/users/{username}` | Modifica username |
| `DELETE` | `/users/{username}` | Cancellazione account |
| `GET` | `/users/{username}/followers` / `/following` | Liste social |
| `PUT`/`DELETE` | `/users/{username}/followers` | Follow / Unfollow |
| `GET`/`PUT`/`DELETE` | `/users/{username}/ban` | Lista ban / Ban / Unban |
| `POST` | `/users/{username}/posts` | Upload di un nuovo post |
| `GET` | `/users/{username}/posts/{postId}` | Recupero immagine binaria del post (pubblico) |
| `DELETE` | `/users/{username}/posts/{postId}` | Cancellazione post |
| `PUT`/`DELETE`/`GET` | `/users/{username}/posts/{postId}/likes` | Like / Unlike / Lista like |
| `POST`/`GET` | `/users/{username}/posts/{postId}/comments` | Aggiunta / Lista commenti |
| `DELETE` | `/users/{username}/posts/{postId}/comments/{commentId}` | Cancellazione commento |
| `GET` | `/stream` | Feed paginato (following + suggeriti) |
| `GET` | `/liveness` | Health check |

Salvo il login, la registrazione e il recupero dell'immagine binaria di un post, tutti gli endpoint richiedono l'header:

```
Authorization: Bearer <sessionToken>
```

## Avvio del progetto

### Backend

Requisiti: Go 1.25+

```sh
# Scarica le dipendenze
go get ./...

# Avvia il server (porta 3000 di default)
go run ./cmd/webapi
```

Il server crea automaticamente il database SQLite e lo schema al primo avvio.

### Frontend

Requisiti: Node.js e npm (oppure Docker, vedi sotto)

```sh
cd webui
npm install
npm run dev
```

Il frontend in modalità sviluppo si aspetta il backend disponibile su `http://localhost:3000` (configurato in `webui/vite.config.js`).

In alternativa, senza installare Node.js localmente, è possibile usare lo script di comodo che apre una shell in un container Node con la cartella del progetto montata:

```sh
./open-npm.sh
```

### Docker

Il backend può essere costruito ed eseguito come immagine Docker multi-stage (build su `golang:1.25-bookworm`, runtime su `debian:bookworm-slim`), con health check integrato:

```sh
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .
docker run -it --rm -p 3000:3000 wasaphoto-backend:latest
```

## Configurazione

Il backend legge la configurazione, in ordine di precedenza crescente, da: variabili d'ambiente (prefisso `CFG`), flag da riga di comando, e infine da un file YAML (che sovrascrive tutto il resto). Un esempio è disponibile in `demo/config.yml`.

Parametri principali:

| Parametro | Default | Descrizione |
|-----------|---------|--------------|
| `Web.APIHost` | `0.0.0.0:3000` | Indirizzo di ascolto del server API |
| `Web.DebugHost` | `0.0.0.0:4000` | Indirizzo del server di debug/profiling |
| `Web.ReadTimeout` / `WriteTimeout` | `5s` | Timeout HTTP |
| `Web.ShutdownTimeout` | `5s` | Timeout per lo shutdown graceful |
| `DB.Filename` | `/tmp/WasaPhoto.db` | Percorso del file SQLite |
| `Config.Path` | `/conf/config.yml` | Percorso del file di configurazione YAML |

## Licenza

Distribuito con licenza [MIT](LICENSE).
