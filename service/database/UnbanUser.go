package database

func (db *appdbimpl) UnbanUser(id string, to_del_id string) error {
	table_name := "\"" + id + "_bans" + "\""
	query1 := "DELETE FROM " + table_name + " WHERE id = ?"
	_, err := db.c.Exec(query1, to_del_id)
	if err != nil {
		return err
	} else {
		return nil
	}

}
