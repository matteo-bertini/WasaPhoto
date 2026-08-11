package database

import "WasaPhoto/service/models"

func (db *appdbimpl) GetStream(requesterId string, limit int, offset int) ([]models.StreamPost, error) {
	var stream []models.StreamPost

	/*
	   Query Logic:
	   1. Select post data and author username by joining 'posts' and 'accounts'.
	   2. Calculate likesNumber and commentsNumber dynamically using subqueries.
	   3. Check 'isLikedByMe' by verifying the existence of a record in the 'likes' table.
	   4. Determine 'isSuggested': 0 if the author is followed, 1 otherwise.
	   5. Filter out posts from banned users (bi-directional check).
	   6. Prioritize followed users' posts, then sort by most recent.
	*/
	query := `
		SELECT 
			p.post_id, 
			a.username, 
			p.caption, 
			p.created_at,
			(SELECT COUNT(*) FROM likes l WHERE l.post_id = p.post_id) as likesNumber,
			(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.post_id) as commentsNumber,
			EXISTS(SELECT 1 FROM likes l WHERE l.post_id = p.post_id AND l.user_id = ?) as isLikedByMe,
			(p.author_id NOT IN (SELECT followed_id FROM follows WHERE follower_id = ?)) as isSuggested
		FROM posts p
		JOIN accounts a ON p.author_id = a.user_id
		WHERE 
			/* Exclude banned users or users who banned the requester */
			p.author_id NOT IN (SELECT banned_id FROM bans WHERE banner_id = ?)
			AND p.author_id NOT IN (SELECT banner_id FROM bans WHERE banned_id = ?)
			
			/* Hybrid Feed: Following + Suggestions (excluding self) */
			AND (
				p.author_id IN (SELECT followed_id FROM follows WHERE follower_id = ?)
				OR 
				(p.author_id NOT IN (SELECT followed_id FROM follows WHERE follower_id = ?) AND p.author_id != ?)
			)
		ORDER BY isSuggested ASC, p.created_at DESC
		LIMIT ? OFFSET ?;
	`

	rows, err := db.c.Query(query,
		requesterId, // for isLikedByMe
		requesterId, // for isSuggested check
		requesterId, // filter bans (I banned them)
		requesterId, // filter bans (They banned me)
		requesterId, // following group
		requesterId, // suggested group
		requesterId, // exclude self from suggested
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.StreamPost
		// Scan row into the models.StreamPost struct
		err := rows.Scan(
			&p.PostId,
			&p.Username,
			&p.Caption,
			&p.DateOfUpload,
			&p.LikesNumber,
			&p.CommentsNumber,
			&p.IsLikedByMe,
			&p.IsSuggested,
		)
		if err != nil {
			return nil, err
		}
		stream = append(stream, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Ensure we return an empty slice instead of nil for consistent JSON
	// encoding: the OpenAPI spec declares the response as `type: array,
	// minItems: 0`, so an empty stream must be encoded as [] and not null.
	if stream == nil {
		stream = []models.StreamPost{}
	}

	return stream, nil
}
