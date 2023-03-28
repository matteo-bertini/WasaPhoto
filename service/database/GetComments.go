package database

func (db *appdbimpl) GetComments(photoid string) (*[]Database_comment, error) {
	table_name := "\"" + photoid + "_comments" + "\""
	query0 := "SELECT * FROM " + table_name
	rows, err := db.c.Query(query0)
	// Si è verificato un errore nell'esecuzione della query
	if err != nil {
		return nil, err

	} else {
		var CommentId string
		var CommentAuthor string
		var CommentText string
		comments := []Database_comment{}
		for rows.Next() {
			err = rows.Scan(&CommentId, &CommentAuthor, &CommentText)
			// Si è verificato un errore nella scan
			if err != nil {
				return nil, err
			} else {
				username, err := db.UsernameFromId(CommentAuthor)
				if err != nil {
					return nil, err
				} else {
					// Username estratto correttamente dal likeid
					comment := Database_comment{CommentId: CommentId, CommentAuthor: *username, CommentText: CommentText}
					comments = append(comments, comment)

				}
			}

		}
		if rows.Err() != nil {
			// Si è verificato un errore durante l'iterazione delle righe o nella loro chiusura
			return nil, rows.Err()
		} else {
			return &comments, nil
		}

	}

}
