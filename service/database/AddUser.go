package database

import "WasaPhoto/service/utils"

func (db *appdbimpl) AddUser(username string, id string) error {
	query1 := "SELECT username FROM users WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		found := rows.Next()
		// Non è stata trovata una entry
		if found == false {
			// Si è verificato un errore nel preparare il risultato o nella chiusura delle righe
			if rows.Err() != nil {
				return err
			} else {
				query2 := "INSERT INTO users VALUES (?,0,0,0)"
				_, err = db.c.Exec(query2, username)
				if err != nil {
					return err
				} else {
					// Creazione delle tabelle photos,followers,following e bans per l'utente (tramite id,in modo che il cambio di username non abbia effetto sul resto)
					table_name := "\"" + id + "_followers" + "\""
					sqlStmt := "CREATE TABLE " + table_name + " (id TEXT NOT NULL PRIMARY KEY)"
					_, err = db.c.Exec(sqlStmt)
					// Errore nella creazione della tabella followers
					if err != nil {
						return err
					} else {
						table_name = "\"" + id + "_following" + "\""
						sqlStmt := "CREATE TABLE " + table_name + " (id TEXT NOT NULL PRIMARY KEY)"
						_, err = db.c.Exec(sqlStmt)
						// Errore nella creazione della tabella following
						if err != nil {
							return err
						} else {
							table_name = "\"" + id + "_bans" + "\""
							sqlStmt := "CREATE TABLE " + table_name + " (id TEXT NOT NULL PRIMARY KEY)"
							_, err = db.c.Exec(sqlStmt)
							// Errore nella creazione della tabella bans
							if err != nil {
								return err

							} else {
								table_name = "\"" + id + "_photos" + "\""
								sqlStmt := "CREATE TABLE " + table_name + " (photoid TEXT NOT NULL PRIMARY KEY,likesnumber INTEGER NOT NULL,commentsnumber INTEGER NOT NULL,dateofupload TEXT NOT NULL)"
								_, err = db.c.Exec(sqlStmt)
								// Errore nella creazione della tabella photos
								if err != nil {
									return err
								} else {
									return nil
								}
							}

						}

					}

				}
			}
		} else {
			err = rows.Close()
			if err != nil {
				return err
			} else {
				return utils.ErrorUserAlreadyExists
			}
		}
	}

}
