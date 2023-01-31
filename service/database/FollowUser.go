package database

import (
	"WasaPhoto/service/utils"
	"errors"
	"strings"
)

func UpdateStringFollow(followers_string string, to_add string) (int, string) {

	// L'utente non ha followers
	if followers_string == "NO" {
		to_ret := to_add
		return 1, to_ret
	} else {
		followers_list := strings.Split(followers_string, ",")
		if utils.CheckPresence(followers_list, to_add) == true {
			return len(followers_list), followers_string
		} else {
			followers_string = followers_string + "," + to_add
			return len(followers_list) + 1, followers_string
		}
	}
}

func (db *appdbimpl) FollowUser(username string, to_add string) error {
	// Un'utente non deve seguire se stesso
	if username == to_add {
		return nil

	}

	// Estrazione della lista dei followers dell'utente dal database
	query1 := "SELECT followers FROM following_followers WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	if err != nil {
		// C'è stato un errore nell'esecuzione della query
		return err
	} else {
		// La query è stata eseguita con successo e ha ritornato un risultato
		if rows.Next() == true {
			var followers_string string
			err = rows.Scan(&followers_string)
			if err != nil {
				return err
			} else {
				err = rows.Close()
				if err != nil {
					return err
				} else {
					// Aggiunta del nuovo follower alla stringa dei followers
					var num_followers int
					num_followers, followers_string = UpdateStringFollow(followers_string, to_add)
					query2 := "UPDATE following_followers SET followers = ? WHERE username = ?"
					_, err = db.c.Exec(query2, followers_string, username)
					if err != nil {
						return err
					} else {
						// Incremento del numero di followers
						query3 := "UPDATE users SET followers = ? WHERE username = ?"
						_, err = db.c.Exec(query3, num_followers, username)
						if err != nil {
							return err
						} else {
							// Estrazione following dalla lista dei following
							query4 := "SELECT following FROM following_followers WHERE username = ?"
							rows, err = db.c.Query(query4, to_add)
							if err != nil {
								// Errore nell'esecuzione della query
								return err
							} else {
								// La query è stata eseguita con successo
								if rows.Next() == true {
									var following_string string
									err = rows.Scan(&following_string)
									if err != nil {
										return err
									} else {
										err = rows.Close()
										if err != nil {
											return err
										} else {
											// Aggiunta del nuovo following alla lista dei following
											var following_number int
											following_number, following_string = UpdateStringFollow(following_string, username)
											query5 := "UPDATE following_followers SET following = ? WHERE username = ?"
											_, err = db.c.Exec(query5, following_string, to_add)
											if err != nil {
												return err
											} else {
												// Incremento del numero di following
												query6 := "UPDATE users SET following = ? WHERE username = ?"
												_, err = db.c.Exec(query6, following_number, to_add)
												if err != nil {
													return err
												} else {
													return nil
												}
											}
										}
									}
								} else {
									if rows.Err() != nil {
										return rows.Err()
									} else {
										err_no_rows := errors.New("La query non ha dato nessun risultato,l'usename non esiste nel database.\n")
										return err_no_rows
									}
								}
							}

						}
					}

				}
			}
		} else {
			if rows.Err() != nil {
				return rows.Err()
			} else {
				err_no_rows := errors.New("La query non ha dato nessun risultato,l'usename non esiste nel database.\n")
				return err_no_rows
			}
		}
	}

}
