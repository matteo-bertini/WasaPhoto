package database

func (db *appdbimpl) DeleteUser(id string, username string) error {
	// Eliminazione dalla tabella users
	/*query1 := "DELETE FROM users WHERE username = ?"
	_, err := db.c.Exec(query1, username)
	if err != nil {
		return err
	} else {
		// Eliminazione dalla tabella authstrings
		query2 := "DELETE FROM authstrings WHERE id = ?"
		_, err := db.c.Exec(query2, id)
		if err != nil {
			return err
		} else{
			// Eliminazione di tutti gli user dalla tabella following dell'user da eliminare
			table_name := "\"" + id + "_following" + "\""
			query3 := "SELECT * FROM " + table_name
			rows, err := db.c.Query(query3)
			if err != nil {
				return err
			} else {
				for rows.Next() == true {
					var to_del_id string
					err = rows.Scan(&to_del_id)
					if err != nil {
						return err
					} else {
						to_del_username, err := db.UsernameFromId(to_del_id)
						if err != nil {
							return err
						} else {
							err = db.UnfollowUser(username, id, *to_del_username, to_del_id)
							if err != nil {
								return nil
							}

						}

					}

				}
				// Si è verificato un errore nell'ottenere il risultato o nel chiudere le rows
				if rows.Err() != nil {
					return rows.Err()
				} else {
					// Rimozione dell'user dalle tabelle dei followers degli user che segue
					table_name = "\"" + id + "_followers" + "\""
					query4 := "SELECT * FROM " + table_name
					rows, err = db.c.Query(query4)
					if err != nil {
						return err
					} else {
						for rows.Next() == true {
							var who_dels_id string
							err = rows.Scan(&who_dels_id)
							if err != nil {
								return err
							} else {
								who_dels_username, err := db.UsernameFromId(who_dels_id)
								if err != nil {
									return err
								} else {
									err = db.UnfollowUser(*who_dels_username, who_dels_id, username, id)
									if err != nil {
										return err
									}

								}

							}
						}
						// Si è verificato un errore nell'ottenere il risultato o nel chiudere le rows
						if rows.Err() != nil {
							return rows.Err()
						} else {
							// Rimozione dei ban
							table_name = "\"" + id + "_bans" + "\""
							query7 := "SELECT * FROM " + table_name
							rows, err = db.c.Query(query7)
							if err != nil {
								return err
							} else {
								for rows.Next() == true {
									var banned_id string
									err = rows.Scan(&banned_id)
									if err != nil {
										return err
									} else {
										banned_username, err := db.UsernameFromId(banned_id)
										if err != nil {
											return err
										} else {
											err = db.UnbanUser(banned_id, *banned_username)
											if err != nil {
												return err
											}
										}

									}
								}
								if rows.Err() != nil {
									return rows.Err()
								} else {
									query8 := "SELECT username FROM users"
									rows, err = db.c.Query(query8)
									if err != nil {
										return err
									} else {
										for rows.Next() == true {
											var username1 string
											err = rows.Scan(&username1)
											if err != nil {
												return err
											} else {
												id1, err := db.IdFromUsername(username1)
												if err != nil {
													return err
												} else {
													err = db.UnbanUser(*id1, id)
													if err != nil {
														return err
													}

												}
											}

										}
										if rows.Err() != nil {
											return rows.Err()
										} else {*/
	table_name := "\"" + id + "_following" + "\""
	query5 := "DROP TABLE " + table_name
	_, err := db.c.Exec(query5)
	if err != nil {
		return err
	} else {
		table_name = "\"" + id + "_followers" + "\""
		query6 := "DROP TABLE " + table_name
		_, err = db.c.Exec(query6)
		if err != nil {
			return err
		} else {
			table_name = "\"" + id + "_bans" + "\""
			query6 := "DROP TABLE " + table_name
			_, err = db.c.Exec(query6)
			if err != nil {
				return err

			} else {
				return nil
			}
		}
	}

	//}
	//}
	//}
	//}
	//}

	//}

	//	}
	//}

	//}
	//}
}
