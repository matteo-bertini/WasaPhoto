package database

func (db *appdbimpl) GetMyStream(id string) (*[]Database_photostream_component, error) {
	// Devo selezionare la foto più recente di ogni user seguito dall'user il cui id è passato come argomento

	table_name := "\"" + id + "_following" + "\""
	query0 := "SELECT id FROM " + table_name
	rows0, err := db.c.Query(query0)
	if err != nil {
		return nil, err
	} else {
		var following_id string
		var photostream_component Database_photostream_component
		stream := []Database_photostream_component{}
		for rows0.Next() == true {
			// Memorizzo l'id del following
			err = rows0.Scan(&following_id)
			if err != nil {
				return nil, err
			} else {
				// Estraggo la foto più recente del following
				table_name = "\"" + following_id + "_photos" + "\""
				query1 := " SELECT *  FROM " + table_name + " ORDER BY dateofupload DESC LIMIT 1"
				rows1, err := db.c.Query(query1)
				if err != nil {
					return nil, err

				} else {
					found := rows1.Next()
					if found == false {
						if rows1.Err() != nil {
							return nil, rows1.Err()
						}
					} else {
						err = rows1.Scan(&photostream_component.PhotoId, &photostream_component.LikesNumber, &photostream_component.CommentsNumber, &photostream_component.DateOfUpload)
						if err != nil {
							return nil, err
						} else {
							following_username, err := db.UsernameFromId(following_id)
							if err != nil {
								return nil, err
							}
							photostream_component.Username = *following_username
							stream = append(stream, photostream_component)
							err = rows1.Close()
							if err != nil {
								return nil, err
							}
						}

					}
				}

			}

		}
		return &stream, nil

	}
}
