package database

func (db *appdbimpl) LikePhoto(userid string, photoid string, likeid string) error {
	table_name := "\"" + photoid + "_likes" + "\""
	query1 := "INSERT INTO " + table_name + " VALUES (?)"
	_, err := db.c.Exec(query1, likeid)
	if err != nil {
		return err
	} else {
		// Aggiorno il numero di likes della foto
		table_name := "\"" + userid + "_photos" + "\""
		query1 := "UPDATE " + table_name + " SET likesnumber=likesnumber+1 WHERE photoid = ?"
		_, err = db.c.Exec(query1, photoid)
		if err != nil {
			return err
		} else {
			return nil
		}

	}

}
