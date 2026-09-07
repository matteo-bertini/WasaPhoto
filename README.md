# WasaPhoto

WasaPhoto è un social network per la condivisione di foto, in stile Instagram. Nasce come progetto per il corso di Web and Software Architectures e viene poi ripreso ed esteso come tesi di laurea triennale in Ingegneria Informatica presso la Sapienza Università di Roma.

Il progetto è composto da un backend REST scritto in Go, un frontend Vue.js 3 e uno storage persistente su SQLite. Le API sono specificate tramite OpenAPI 3.0 (`doc/api.yaml`).

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

- **Autenticazione**: login e registrazione condivisi su un unico endpoint (`isSignUp`), sessione gestita con token Bearer.
- **Profilo**: visualizzazione con contatori (follower, following, post), modifica dello username, cancellazione dell'account.
- **Post**: upload di una foto con didascalia opzionale, cancellazione, recupero dell'immagine binaria.
- **Interazioni**: like/unlike sui post, commenti (creazione, lettura, cancellazione).
- **Relazioni sociali**: follow/unfollow, liste di follower e following.
- **Ban**: un utente può bloccarne un altro; il ban nasconde reciprocamente profilo, post, stream, like e liste social, ed elimina automaticamente un eventuale follow esistente.
- **Stream**: feed paginato che dà priorità ai post degli utenti seguiti e completa gli slot rimanenti con post suggeriti, escludendo sempre gli utenti bannati.

## Stack tecnologico

**Backend**
- [Go](https://go.dev/) 1.25
- [httprouter](https://github.com/julienschmidt/httprouter) per il routing
- [go-sqlite3](https://github.com/mattn/go-sqlite3) come driver SQLite
- [logrus](https://github.com/sirupsen/logrus) per il logging strutturato
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) per l'hashing delle password
- [ksuid](https://github.com/segmentio/ksuid) per identificatori univoci e ordinabili (utenti, post, sessioni)
- [ardanlabs/conf](https://github.com/ardanlabs/conf) per la configurazione da flag, env e YAML
- [gorilla/handlers](https://github.com/gorilla/handlers) per il middleware CORS

**Frontend**
- [Vue.js 3](https://vuejs.org/)
- [Vue Router 4](https://router.vuejs.org/)
- [Axios](https://axios-http.com/)
- [Vite](https://vitejs.dev/) come build tool e dev server

**Storage**
- SQLite, con schema creato automaticamente all'avvio e vincoli `FOREIGN KEY ... ON DELETE CASCADE`: cancellare un account rimuove a cascata profilo, post, like, commenti, follow e ban collegati.

## Struttura del progetto

```
WasaPhoto/
├── cmd/
│   ├── webapi/            # entry point del server: main, configurazione, CORS, embedding della web UI
│   └── healthcheck/       # binario usato dal probe HEALTHCHECK del container Docker
├── service/
│   ├── api/               # handler HTTP, routing, middleware di autenticazione
│   ├── database/          # accesso ai dati: schema e query
│   └── models/            # strutture dati condivise ed errori di dominio
├── webui/                 # applicazione frontend Vue.js
├── doc/
│   └── api.yaml           # specifica OpenAPI 3.0 delle API REST
├── demo/
│   └── config.yml         # esempio di file di configurazione
├── Dockerfile.backend     # build multi-stage del backend Go
├── Dockerfile.frontend    # build multi-stage del frontend (Node + Nginx)
└── open-npm.sh            # apre una shell in un container Node per lavorare sul frontend senza installare nulla in locale
```

## Architettura

Il backend è organizzato su tre livelli con responsabilità separate:

- **`service/api`**: espone gli endpoint REST, valida gli input e traduce gli errori di dominio in codici di stato HTTP. L'autenticazione è centralizzata in un middleware (`AuthMiddleware`) che verifica il token Bearer e inietta lo `UserID` nel contesto della richiesta.
- **`service/database`**: nasconde tutta la logica SQL dietro l'interfaccia `AppDatabase`, così l'implementazione concreta (SQLite) resta disaccoppiata dagli handler.
- **`service/models`**: definisce le strutture condivise (Post, Comment, UserProfile, ecc.) e gli errori di dominio (`ErrUserNotFound`, `ErrForbiddenAction`, ...) usati per mappare le risposte HTTP.

Le operazioni che toccano più tabelle e devono restare atomiche — registrazione, like con verifica del ban, follow, aggiunta di un commento — vengono eseguite dentro transazioni SQL.

Il frontend è una SPA Vue.js che parla con il backend via REST e mantiene la sessione (token e username) in `localStorage`.

Per il deployment sono previste due modalità, entrambe supportate dal codice:

- **due container separati**, backend Go e frontend Nginx, che è il caso d'uso pensato per `Dockerfile.backend` e `Dockerfile.frontend`;
- **binario singolo**, compilando con il build tag `webui` (`go build -tags webui ./cmd/webapi`): in questo caso i file statici del frontend vengono incorporati nell'eseguibile Go tramite `go:embed` e serviti sotto `/dashboard/`, mentre le API restano sulla root. Senza il tag, `registerWebUI` è uno stub e il binario serve solo le API.

## Modello dati

Lo schema SQLite viene creato automaticamente all'avvio (`service/database/database.go`):

| Tabella     | Descrizione                                      |
|-------------|---------------------------------------------------|
| `accounts`  | Credenziali e token di sessione                   |
| `profiles`  | Bio e avatar dell'utente                          |
| `posts`     | Post pubblicati (didascalia, autore, data)        |
| `follows`   | Relazioni di follow (follower → followed)         |
| `likes`     | Like sui post                                     |
| `comments`  | Commenti sui post                                 |
| `bans`      | Relazioni di ban (banner → banned)                |

Tutte le tabelle sono legate ad `accounts` con `ON DELETE CASCADE`: cancellare un account rimuove automaticamente profilo, post, like, commenti, follow e ban associati.

## API REST

Le API sono documentate per intero in `doc/api.yaml` (OpenAPI 3.0). Endpoint principali:

| Metodo | Path | Descrizione |
|--------|------|-------------|
| `POST` | `/session` | Login o registrazione (`isSignUp`) |
| `DELETE` | `/session` | Logout |
| `GET` | `/users/{username}` | Profilo utente, statistiche e post |
| `PUT` | `/users/{username}` | Modifica username |
| `DELETE` | `/users/{username}` | Cancellazione account |
| `GET` | `/users/{username}/followers` \| `/following` | Liste social |
| `PUT` / `DELETE` | `/users/{username}/followers` | Follow / Unfollow |
| `GET` / `PUT` / `DELETE` | `/users/{username}/ban` | Lista ban / Ban / Unban |
| `POST` | `/users/{username}/posts` | Upload di un nuovo post |
| `GET` | `/users/{username}/posts/{postId}` | Immagine binaria del post (endpoint pubblico) |
| `DELETE` | `/users/{username}/posts/{postId}` | Cancellazione post |
| `PUT` / `DELETE` / `GET` | `/users/{username}/posts/{postId}/likes` | Like / Unlike / Lista like |
| `POST` / `GET` | `/users/{username}/posts/{postId}/comments` | Aggiunta / Lista commenti |
| `DELETE` | `/users/{username}/posts/{postId}/comments/{commentId}` | Cancellazione commento |
| `GET` | `/stream` | Feed paginato (following + suggeriti) |
| `GET` | `/liveness` | Health check, usato dal probe Docker |

Ad eccezione di login, registrazione e recupero dell'immagine di un post, tutti gli endpoint richiedono l'header:

```
Authorization: Bearer <sessionToken>
```

## Avvio del progetto

### Backend

Richiede Go 1.25+.

```sh
go get ./...
go run ./cmd/webapi
```

Al primo avvio il server crea da sé il database SQLite e lo schema (percorso di default: `/tmp/WasaPhoto.db`).

### Frontend

Richiede Node.js e npm (in alternativa vedi Docker, sotto).

```sh
cd webui
npm install
npm run dev
```

In modalità sviluppo il frontend si aspetta il backend su `http://localhost:3000` (impostato in `webui/vite.config.js`).

Se non si vuole installare Node.js in locale, è disponibile uno script che apre una shell in un container Node con la cartella del progetto montata:

```sh
./open-npm.sh
```

### Docker

Backend e frontend sono pensati come due immagini separate.

Backend (Go, build multi-stage su `golang:1.25-bookworm` e runtime su `debian:bookworm-slim`, con health check integrato che interroga `/liveness`):

```sh
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .
docker run -it --rm -p 3000:3000 wasaphoto-backend:latest
```

Frontend (build con `node:24-bookworm-slim`, servito poi da `nginx:1.27-alpine`):

```sh
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .
docker run -it --rm -p 8081:80 wasaphoto-frontend:latest
```

## Configurazione

Il backend legge la configurazione da tre fonti, in ordine di precedenza crescente: variabili d'ambiente (prefisso `CFG`), flag da riga di comando, e infine un file YAML, che sovrascrive tutto il resto. Un esempio è in `demo/config.yml`.

| Parametro | Default | Descrizione |
|-----------|---------|--------------|
| `Web.APIHost` | `0.0.0.0:3000` | Indirizzo di ascolto dell'API |
| `Web.DebugHost` | `0.0.0.0:4000` | Indirizzo del server di debug/profiling |
| `Web.ReadTimeout` / `WriteTimeout` | `5s` | Timeout HTTP |
| `Web.ShutdownTimeout` | `5s` | Timeout per lo shutdown graceful |
| `DB.Filename` | `/tmp/WasaPhoto.db` | Percorso del file SQLite |
| `Config.Path` | `/conf/config.yml` | Percorso del file di configurazione YAML |

## Licenza

Distribuito con licenza [MIT](LICENSE).
