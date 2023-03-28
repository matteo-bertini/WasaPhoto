package database

// Questa funzione ausiliaria rimuove correttamente tutti i followers dell'user passato come argomento
func (db *appdbimpl) RemoveAllFollowers(id string, username string) error {
	table_name := "\"" + id + "_followers" + "\""
	query1 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query1)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		// La query è stata eseguita correttamente
		var follower_id string
		followers := []string{}
		// Aggiungo tutti i followers alla lista per toglierli dopo
		for rows.Next() {
			err = rows.Scan(&follower_id)
			// Si è verificato un errore nello scan
			if err != nil {
				return err

			} else {
				// Aggiungo il follower_id alla lista
				followers = append(followers, follower_id)
			}

		}
		// Si è verificato un errore nella preparazione del risultato o nella chiusura automatica delle rows
		if rows.Err() != nil {
			return rows.Err()

		} else {
			// Le rows sono state chiuse automaticamente
			for _, f := range followers {
				follower_username, err := db.UsernameFromId(f)
				// Si è verificato un errore
				if err != nil {
					return err
				} else {
					err = db.UnfollowUser(username, id, *follower_username, f)
					// Si è verificato un errore
					if err != nil {
						return err
					}
				}

			}
			return nil

		}

	}
}

// Questa funzione ausiliaria rimuove correttamente tutti i following dell'user passato come argomento
func (db *appdbimpl) RemoveAllFollowing(id string, username string) error {
	table_name := "\"" + id + "_following" + "\""
	query1 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query1)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		// La query è stata eseguita correttamente
		var following_id string
		following := []string{}
		// Aggiungo tutti i following alla lista per toglierli dopo
		for rows.Next() {
			// Si è verificato un errore nello scan
			err = rows.Scan(&following_id)
			if err != nil {
				return err

			} else {
				following = append(following, following_id)
			}

		}
		// Si è verificato un errore nella preparazione del risultato o nella chiusura automatica delle rows
		if rows.Err() != nil {
			return rows.Err()

		} else {
			// Le rows sono state chiuse automaticamente
			for _, f := range following {
				following_username, err := db.UsernameFromId(f)
				// Si è verificato un errore
				if err != nil {
					return err
				} else {
					err = db.UnfollowUser(*following_username, f, username, id)
					// Si è verificato un errore
					if err != nil {
						return err
					}
				}

			}
			return nil

		}

	}
}

// Questa funzione ausiliaria rimuove correttamente tutte le foto (e le tabelle ad esse correlate) dell'user passato come argomento
func (db *appdbimpl) RemoveAllPhotos(id string, username string) error {
	table_name := "\"" + id + "_photos" + "\""
	query1 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query1)

	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return err
	} else {
		// La query è stata eseguita correttamente
		var photo Database_photo
		photos := []string{}
		// Aggiungo tutte le foto alla lista per toglierle dopo
		for rows.Next() {
			err = rows.Scan(&photo.PhotoId, &photo.LikesNumber, &photo.CommentsNumber, &photo.DateOfUpload)
			// Si è verificato un errore nello scan
			if err != nil {
				return err

			} else {
				photos = append(photos, photo.PhotoId)
			}

		}
		// Si è verificato un errore nella preparazione del risultato o nella chiusura automatica delle rows
		if rows.Err() != nil {
			return rows.Err()

		} else {
			// Le rows sono state chiuse automaticamente
			for _, ph := range photos {
				err = db.DeletePhoto(id, ph)
				// Si è verificato un errore
				if err != nil {
					return err
				}
			}
			return nil

		}

	}
}

// Questa funzione rimmuove tutte le tabelle associate all'user
func (db *appdbimpl) RemoveAllTables(id string, username string) error {
	// Eliminazione della tabella followers
	table_name := "\"" + id + "_followers" + "\""
	stmt1 := "DROP TABLE " + table_name
	_, err := db.c.Exec(stmt1)
	if err != nil {
		return err
	} else {
		// Eliminazione della tabella following
		table_name = "\"" + id + "_following" + "\""
		stmt1 := "DROP TABLE " + table_name
		_, err := db.c.Exec(stmt1)
		if err != nil {
			return err
		} else {
			// Eliminazione della tabella bans
			table_name = "\"" + id + "_bans" + "\""
			stmt1 := "DROP TABLE " + table_name
			_, err := db.c.Exec(stmt1)
			if err != nil {
				return err
			} else {
				// Eliminazione della tabella photos
				table_name = "\"" + id + "_photos" + "\""
				stmt1 := "DROP TABLE " + table_name
				_, err := db.c.Exec(stmt1)
				if err != nil {
					return err
				} else {
					return nil

				}

			}

		}

	}
}
func (db *appdbimpl) DeleteUser(id string, username string) error {
	// Rimuovo tutti i followers
	err := db.RemoveAllFollowers(id, username)
	if err != nil {
		return err
	} else {
		// Rimuovo tutti i following
		err = db.RemoveAllFollowing(id, username)
		if err != nil {
			return err
		} else {
			// Rimuovo tutte le foto e le tabelle correlate alle foto
			err = db.RemoveAllPhotos(id, username)
			if err != nil {
				return err
			} else {
				// Rimuovo l'username dalla tabella users
				stmt := "DELETE FROM users WHERE username = ?"
				_, err = db.c.Exec(stmt, username)
				if err != nil {
					return err
				} else {
					// Rimuovo l'id dalla tabella authstrings
					stmt := "DELETE FROM authstrings WHERE id = ?"
					_, err = db.c.Exec(stmt, id)
					if err != nil {
						return err
					} else {
						// Rimuovo tutte le tables relative all'user
						err = db.RemoveAllTables(id, username)
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

}
