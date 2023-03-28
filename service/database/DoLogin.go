package database

import (
	"github.com/segmentio/ksuid"
)

func (db *appdbimpl) DoLogin(username string) (*string, error) {
	// Controllo che l'user con username specificato non sia già registrato

	query1 := "SELECT id FROM authstrings WHERE username = ?"
	rows, err := db.c.Query(query1, username)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		// La query è stata eseguita correttamente, ma non è stato possibile preparare il risultato
		if !rows.Next() {
			err = rows.Err()
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			if err != nil {
				return nil, err
			} else {
				// La query non ha dato nessun risultato,quindi devo creare una nuova entry
				id := ksuid.New().String()
				query2 := "INSERT INTO authstrings VALUES (?,?)"
				_, err := db.c.Exec(query2, username, id)
				if err != nil {
					return nil, err
				} else {
					return &id, nil
				}

			}
			// Esiste una entry nella tabella con l'username specificato quindi estraggo l'Id
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
