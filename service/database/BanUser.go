package database

import (
	"WasaPhoto/service/utils"
	"errors"
	"strings"
)

func addBan(bannedstring string, to_add string) string {
	if bannedstring == "NO" {
		return to_add
	} else {
		banned_list := strings.Split(bannedstring, ",")
		if utils.CheckPresence(banned_list, to_add) == true {
			return bannedstring
		} else {
			bannedstring = bannedstring + "," + to_add
			return bannedstring
		}
	}
}

func (db *appdbimpl) BanUser(username string, to_ban string) error {
	// Estrazione della lista dei bannati di "username"
	query1 := "SELECT has_banned FROM ban_resume WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	defer rows.Close()
	if err != nil {
		return err
	} else {
		if rows.Next() == true {
			var hasbanned string
			err = rows.Scan(&hasbanned)
			if err != nil {
				return err
			} else {
				err = rows.Close()
				if err != nil {
					return err
				} else {
					hasbanned = addBan(hasbanned, to_ban)
					query2 := "UPDATE ban_resume SET has_banned = ? WHERE username = ?"
					_, err = db.c.Exec(query2, hasbanned, username)
					if err != nil {
						return err
					} else {
						query3 := "SELECT banned_by FROM ban_resume WHERE username = ?"
						rows, err = db.c.Query(query3, to_ban)
						if err != nil {
							return err
						} else {
							if rows.Next() == true {
								var bannedby string
								err = rows.Scan(&bannedby)
								if err != nil {
									return err
								} else {
									err = rows.Close()
									if err != nil {
										return err
									} else {
										bannedby = addBan(bannedby, username)
										query4 := "UPDATE ban_resume SET banned_by = ? WHERE username = ?"
										_, err = db.c.Exec(query4, bannedby, to_ban)
										if err != nil {
											return err
										} else {
											return nil
										}
									}
								}
							} else {
								if rows.Err() != nil {
									return rows.Err()
								} else {
									return errors.New("La query non ha dato nessun risultato,l'usename non esiste nel database.\n")
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
				return errors.New("La query non ha dato nessun risultato,l'usename non esiste nel database.\n")
			}
		}
	}

}
