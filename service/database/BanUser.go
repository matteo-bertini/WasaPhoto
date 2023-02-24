package database

import "WasaPhoto/service/utils"

func (db *appdbimpl) BanUser(username string, id string, to_ban_username string, to_ban_id string) error {
	table_name := "\"" + id + "_bans" + "\""
	query0 := "SELECT * FROM " + table_name + " WHERE id = ?"
	rows, err := db.c.Query(query0, to_ban_id)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		// La query è stata eseguita correttamente
		if rows.Next() == false {
			// Si è verificato un errore nel preparare il risultato o nel chiudere le rows
			if rows.Err() != nil {
				return rows.Err()
			} else {
				// Nella tabella bans non esiste l'id cercato
				// Chiudo le righe della tabella bans
				err = rows.Close()
				if err != nil {
					return err
				} else {
					// Inserisco l'id nella tabella bans
					query1 := "INSERT INTO " + table_name + " VALUES (?)"
					_, err := db.c.Exec(query1, to_ban_id)
					// Si è verificato un errore nell'esecuzione dello statement
					if err != nil {
						return err
					} else {
						// Rimuovo il follow reciprocamente dopo il ban
						err = db.UnfollowUser(username, id, to_ban_username, to_ban_id)
						if err != nil {
							return err
						} else {
							err = db.UnfollowUser(to_ban_username, to_ban_id, username, id)
							if err != nil {
								return err
							} else {
								return nil

							}
						}
					}

				}

			}
		} else {
			err = rows.Close()
			if err != nil {
				return err
			}
			return utils.ErrUserAlreadyBanned
		}
	}

}
