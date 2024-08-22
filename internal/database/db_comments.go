package database

import "forumProject/internal/models"

func GetComments(postID int) ([]models.Comment, error) {
	query := `
		SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.likes, comments.dislikes, comments.created_at, users.username
		FROM comments
		JOIN users ON comments.user_id = users.id
		WHERE comments.post_id = ?
	`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var username string
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.CreatedAt, &username)
		if err != nil {
			return nil, err
		}
		comment.Username = username
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func CreateComment(comment models.Comment) error {
	sqlStm := `INSERT INTO comments (
	post_id, user_id, 
	content, created_at) 
	VALUES (?, ?, ?, strftime('%Y-%m-%d %H:%M:%S', 'now', '+3 hours'))
	`

	stmt, err := DB.Prepare(sqlStm)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}

	// commentID, err := result.LastInsertId()
	// if err != nil {
	// 	return err
	// }

	return nil
}
