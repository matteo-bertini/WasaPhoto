package database

import (
	"WasaPhoto/service/utils"
	"errors"
	"strings"
)

func UpdateBanString(banstring string, to_del string) string {
	if banstring == "NO" {
		return banstring
	} else {
		banlist := strings.Split(banstring, ",")
		if utils.CheckPresence(banlist, to_del) == false {
			return banstring
		} else {
			var new_list []string
			for _, el := range banlist {
				if el != to_del {
					new_list = append(new_list, el)
				}
			}
			to_ret := strings.Join(new_list, ",")
			if to_ret == "" {
				return "NO"
			} else {
				return to_ret
			}

		}
	}
}
func (db *appdbimpl) UnbanUser(username string, to_del string) error {

	// Controllo l'esistenza di entrambi gli utenti
	err1, ex1 := db.CheckExistence(username)
	err2, ex2 := db.CheckExistence(to_del)
	if err1 != nil || err2 != nil {
		return errors.New("InternalServerError")
	} else {
		if *ex1 == false || *ex2 == false {
			return errors.New("NotFound")
		} else {
			// Entrambi gli utenti esistono nel database

			// Elimino to_del dalla lista bannati di username
			query1 := "SELECT has_banned FROM ban_resume WHERE username = ?"
			rows, err := db.c.Query(query1, username)
			if err != nil {
				return err
			} else {
				// Siccome entrambi gli utenti esistono nel database,per come è inizializzato (alla creazione del profilo inizializzo le tabelle)
				// non c'è bisogno di controllare rows.Next() == true,sono sicuro che ci sarà un riga! (al massimo settata tutta  a "NO")
				rows.Next()
				var has_banned_string string
				err = rows.Scan(&has_banned_string)
				if err != nil {
					return err
				} else {
					err = rows.Close()
					if err != nil {
						return err
					} else {
						has_banned_string = UpdateBanString(has_banned_string, to_del)
						query2 := "UPDATE ban_resume SET has_banned = ? WHERE username= ?"
						_, err = db.c.Exec(query2, has_banned_string, username)
						if err != nil {
							return err
						} else {
							query3 := "SELECT banned_by FROM ban_resume WHERE username = ?"
							rows, err := db.c.Query(query3, to_del)
							if err != nil {
								return err
							} else {
								rows.Next()
								var banned_by_string string
								err = rows.Scan(&banned_by_string)
								if err != nil {
									return err
								} else {
									err = rows.Close()
									if err != nil {
										return err
									} else {
										banned_by_string = UpdateBanString(banned_by_string, username)
										query4 := "UPDATE ban_resume SET banned_by = ? WHERE username= ?"
										_, err = db.c.Exec(query4, banned_by_string, to_del)
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
			}

		}
	}

}
