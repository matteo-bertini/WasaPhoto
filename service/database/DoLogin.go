package database

import (
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

func (db *appdbimpl) DoLogin(username string, password string) (*string, error, *bool) {

	// Controllo che l'user con username e password specificati non sia già registrato

	query1 := "SELECT id,hashedpassword FROM authstrings WHERE username = ?"
	rows, err := db.c.Query(query1, username)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err, nil
	} else {
		// La query è stata eseguita correttamente, ma non è stato possibile preparare il risultato
		if !rows.Next() {
			err = rows.Err()
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			if err != nil {
				return nil, err, nil
			} else {
				// La query non ha dato nessun risultato,quindi devo creare una nuova entry
				id := ksuid.New().String()

				// Faccio l'hashing della password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
				// Si è verificato un errore nella generazione dell'hash
				if err != nil {
					return nil, err, nil
				}

				query2 := "INSERT INTO authstrings VALUES (?,?,?)"
				_, err = db.c.Exec(query2, username, string(hashedPassword), id)
				if err != nil {
					return nil, err, nil
				} else {
					var created = true
					return &id, nil, &created
				}

			}

			// Esiste una entry nella tabella con l'username specificato quindi controllo che la password sia corretta ed estraggo l'Id
		} else {
			var id string
			var hashedPassword string

			err = rows.Scan(&id, &hashedPassword)
			if err != nil {
				return nil, err, nil
			} else {
				err = rows.Close()
				if err != nil {
					return nil, err, nil
				} else {
					err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
					if err != nil {
						return nil, err, nil
					}
					var created = false
					return &id, nil, &created
				}
			}

		}
	}

}
