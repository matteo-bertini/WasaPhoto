package database

import "WasaPhoto/service/utils"

func (db *appdbimpl) GetUserProfile(username string) (*Database_user, error) {
	query1 := "SELECT * FROM users WHERE username = ?"
	var db_user Database_user
	rows, err := db.c.Query(query1, username)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err
	} else {
		found := rows.Next()
		// é stato trovato un record (user profile) nella tabella users
		if found {
			// Memorizzo i valori del record trovato nella struct da ritornare
			err = rows.Scan(&db_user.Username, &db_user.Followers, &db_user.Following, &db_user.Numberofphotos)
			if err != nil {
				return nil, err
			} else {
				err = rows.Close()
				if err != nil {
					return nil, err
				} else {
					// Devo memorizzare le photos dell'user nell'ultimo campo della struct di ritorno
					id, err := db.IdFromUsername(username)
					if err != nil {
						return nil, err
					} else {
						table_name := "\"" + *id + "_photos" + "\""
						query2 := "SELECT * FROM " + table_name
						rows, err = db.c.Query(query2)
						if err != nil {
							return nil, err
						} else {
							uploaded_photos := []Database_photo{}
							for rows.Next() {
								var photo Database_photo
								err = rows.Scan(&photo.PhotoId, &photo.LikesNumber, &photo.CommentsNumber, &photo.DateOfUpload)
								if err != nil {
									return nil, err
								} else {
									uploaded_photos = append(uploaded_photos, photo)

								}

							}
							if rows.Err() != nil {
								return nil, rows.Err()
							} else {
								db_user.UploadedPhotos = uploaded_photos
								return &db_user, nil
							}

						}

					}

				}
			}
		} else {
			if rows.Err() != nil {
				return nil, rows.Err()
			} else {
				return nil, utils.ErrUserDoesNotExist
			}
		}

	}
}
