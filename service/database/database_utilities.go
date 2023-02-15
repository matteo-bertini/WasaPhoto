package database

import (
	"WasaPhoto/service/utils"
	"net/http"
	"strings"
)

// La funzione CheckAuthorization controlla che l'Id passato nella richiesta corrisponda all'username passato come argomento
func (db *appdbimpl) CheckAuthorization(request *http.Request, username string) error {
	// Estrazione dell' Authorization dall'header

	auth_header := request.Header.Get("Authorization")

	// Non è stato specificato il campo Authorization nell'header
	if auth_header == "" {
		return utils.ErrorAuthorizationNotSpecified
	} else {
		// Il campo Authorization è stato specificato nell'header (Authorization : Bearer abcdef)
		splitted_authorization := strings.Split(auth_header, " ")

		// Controllo che l'Id sia stato specificato correttamente nel campo Authorization
		if len(splitted_authorization) == 2 {
			authorization_type := splitted_authorization[0]
			id := splitted_authorization[1]
			// Id non specificato in conformità con le specifiche
			if authorization_type != "Bearer" || strings.TrimSpace(id) == "" {
				return utils.ErrorBearerTokenNotSpecifiedWell
			} else {
				// Id specificato nel campo Authorization in modo corretto
				query1 := "SELECT * FROM authstrings WHERE id = ? AND username = ?"
				rows, err := db.c.Query(query1, id, username)
				if err != nil {
					// Si è verificato un errore nell'esecuzione della query
					return err
				} else {
					found := rows.Next()
					if found == false {
						// Si è verififcato un errore nella preparazione della result row o nella chiusura delle rows (se rows.Next()==false vengono chiuse automaticamente)
						if rows.Err() != nil {
							return rows.Err()

						} else {
							// Non è stato trovata una entry nel database con username ed id specificati
							return utils.ErrorUnauthorized

						}

					} else {
						err = rows.Close()
						if err != nil {
							return err
						} else {
							// Autorizzato
							return nil
						}
					}
				}
			}
		} else {
			return utils.ErrorBearerTokenNotSpecifiedWell
		}
	}
}

// La funzione CheckUserExistence controlla che l'username passato nella richiesta sia un utente (registrato e con profilo esistente)
func (db *appdbimpl) CheckUserExistence(username string) error {
	// Cerco una entry nella tabella users con  username passato come argomento
	query1 := "SELECT * FROM users WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		found := rows.Next()
		// Non è stata trovata una entry o c'è stato un errore
		if found == false {
			// Si è verificato un errore nella preparazione del risultato o nella chiusura delle righe
			if rows.Err() != nil {
				return rows.Err()
			} else {
				// Non è stata trovata una entry nel database
				return utils.ErrUserDoesNotExist
			}

		} else {
			// é stata trovata una entry nel database,l'utente esiste
			err = rows.Close()
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}

}

// La funzione IdFromUsername restituisce l'Id dell'username passato
func (db *appdbimpl) IdFromUsername(username string) (*string, error) {
	query1 := "SELECT id FROM authstrings WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		found := rows.Next()
		// Non è stata trovata una entry o c'è stato un errore
		if found == false {
			// Si è verificato un errore nella preparazione del risultato o nella chiusura delle righe
			if rows.Err() != nil {
				return nil, rows.Err()
			} else {
				// Non è stata trovata una entry nel database
				return nil, utils.ErrUserDoesNotExist
			}
		} else {
			var id string
			err = rows.Scan(&id)
			if err != nil {
				return nil, err
			} else {
				err = rows.Close()
				if err != nil {
					return nil, err
				} else {
					return &id, nil
				}
			}

		}

	}
}

// La funzione UsernameFromId restituisce l'Id dell'username passato
func (db *appdbimpl) UsernameFromId(id string) (*string, error) {
	query1 := "SELECT username FROM authstrings WHERE id = ?"
	rows, err := db.c.Query(query1, id)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		found := rows.Next()
		// Non è stata trovata una entry o c'è stato un errore
		if found == false {
			// Si è verificato un errore nella preparazione del risultato o nella chiusura delle righe
			if rows.Err() != nil {
				return nil, rows.Err()
			} else {
				// Non è stata trovata una entry nel database
				return nil, utils.ErrUserDoesNotExist
			}
		} else {
			var username string
			err = rows.Scan(&username)
			if err != nil {
				return nil, err
			} else {
				err = rows.Close()
				if err != nil {
					return nil, err
				} else {
					return &username, nil
				}
			}

		}

	}
}

// La funzione IsAllowed controlla che l'user con id1 sia in autorizzato a vedere le informazioni dell'user con id2
// Questa funzione assume come condizione che entrambi gli utenti esistano
func (db *appdbimpl) IsAllowed(id1 string, id2 string) error {

	// Controllo che id1 non abbia bannato id2
	table_name := "\"" + id1 + "_bans" + "\""
	query1 := "SELECT id FROM " + table_name + " WHERE id = ?"
	rows, err := db.c.Query(query1, id2)
	if err != nil {
		return err
	} else {
		found := rows.Next()
		if found == true {
			return utils.ErrUserNotAllowed
		} else {
			if rows.Err() != nil {
				return rows.Err()
			} else {
				// id1 non ha bannato id2
				table_name := "\"" + id2 + "_bans" + "\""
				query2 := "SELECT id FROM " + table_name + " WHERE id = ?"
				rows, err := db.c.Query(query2, id1)
				if err != nil {
					return err
				} else {
					found := rows.Next()
					if found == true {
						return utils.ErrUserNotAllowed
					} else {
						if rows.Err() != nil {
							return rows.Err()
						} else {
							// id2 non ha bannato id1
							return nil
						}
					}
				}
			}
		}
	}
}

// La funzione CheckPhotoExistence controlla che la foto con l'id passato esista nel profilo dell'user con l'id passato
func (db *appdbimpl) CheckPhotoExistence(user_id string, photoid string) error {
	// Cerco se esiste una foto con l'id passato all'interno della tabella photos dell'user passato
	table_name := "\"" + user_id + "_photos" + "\""
	query0 := "SELECT * FROM " + table_name + " WHERE photoid = ?"
	rows, err := db.c.Query(query0, photoid)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err

	} else {
		// Query eseguita correttamente
		found := rows.Next()
		// Non è stato preparato un risultato
		if found == false {
			// Si è verificato un errore nella preparazione del risultato o nella chiusura delle rows
			if rows.Err() != nil {
				return rows.Err()
			} else {
				// La foto non esiste
				return utils.ErrPhotoDoesNotExist

			}
		} else {
			// La foto esiste
			err = rows.Close()
			if err != nil {
				return err
			} else {
				return nil
			}

		}
	}

}
