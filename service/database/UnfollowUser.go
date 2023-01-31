package database

import (
	"WasaPhoto/service/utils"
	"errors"
	"strings"
)

func UpdateStringUnfollow(fstring string, to_del string) (int, string) {
	if fstring == "NO" {
		return 0, fstring
	} else {
		flist := strings.Split(fstring, ",")
		if utils.CheckPresence(flist, to_del) == false {
			return len(flist), fstring
		} else {
			var new_list []string
			for _, el := range flist {
				if el != to_del {
					new_list = append(new_list, el)
				}
			}
			to_ret := strings.Join(new_list, ",")
			if to_ret == "" {
				return 0, "NO"
			} else {
				return len(new_list), to_ret
			}

		}
	}

}
func (db *appdbimpl) UnfollowUser(username string, to_del string) error {
	if username == to_del {
		return nil
	}
	query1 := "SELECT followers FROM following_followers WHERE username = ?"
	rows, err := db.c.Query(query1, username)
	if err != nil {
		return err
	} else {
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
					var followers_number int
					followers_number, followers_string = UpdateStringUnfollow(followers_string, to_del)
					query2 := "UPDATE following_followers SET followers = ? WHERE username = ?"
					_, err = db.c.Exec(query2, followers_string, username)
					if err != nil {
						return err
					} else {
						// Decremento numero dei followers
						query3 := "UPDATE users SET followers = ? WHERE username = ?"
						_, err = db.c.Exec(query3, followers_number, username)
						if err != nil {
							return err
						} else {
							// Rimozione dell username dalla lista dei following di to_del
							query4 := "SELECT following FROM following_followers WHERE username =?"
							rows, err = db.c.Query(query4, to_del)
							if err != nil {
								return err
							} else {
								if rows.Next() == false {
									if rows.Err() != nil {
										return rows.Err()
									} else {
										return errors.New("The username searched does not exists.")
									}
								}
								var following_string string
								err = rows.Scan(&following_string)
								if err != nil {
									return err
								} else {
									err = rows.Close()
									if err != nil {
										return err
									}

									var following_number int
									following_number, following_string = UpdateStringUnfollow(following_string, username)
									query5 := "UPDATE following_followers SET following = ? WHERE username = ?"
									_, err = db.c.Exec(query5, following_string, to_del)
									if err != nil {
										return err
									} else {
										query6 := "UPDATE users SET following = ? WHERE username = ?"
										_, err = db.c.Exec(query6, following_number, to_del)
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

		} else {
			if rows.Err() != nil {
				return rows.Err()
			} else {
				return errors.New("The username searched does not exists.")
			}
		}
	}

}
