package database

import (
	"database/sql"
	"fmt"
	"forumProject/internal/models"
)

func Liking(entityType, ltype string, id, userID int) ([]int, error) {

	var ltypei int
	if ltype == "like" {
		ltypei = 1
	} else if ltype == "dislike" {
		ltypei = -1
	} else {
		return []int{}, fmt.Errorf("invalid like type: %s", ltype)
	}

	// Check if the user has already liked or disliked the entity
	var existingLikeType int
	var sqlStm string
	if entityType == "post" {
		sqlStm = `SELECT like_type FROM likes WHERE user_id = ? AND post_id = ?`
	} else if entityType == "comment" {
		sqlStm = `SELECT like_type FROM likes WHERE user_id = ? AND comment_id = ?`
	} else {
		return []int{}, fmt.Errorf("invalid entity type: %s", entityType)
	}

	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return []int{}, err
	}
	defer stm.Close()

	err = stm.QueryRow(userID, id).Scan(&existingLikeType)
	if err != nil && err != sql.ErrNoRows {
		return []int{}, err
	}

	if err == sql.ErrNoRows {
		// User has not liked or disliked the entity, insert a new like
		if entityType == "post" {
			sqlStm = `INSERT INTO likes (user_id, post_id, like_type, created_at) VALUES (?, ?, ?, datetime('now'))`
		} else if entityType == "comment" {
			sqlStm = `INSERT INTO likes (user_id, comment_id, like_type, created_at) VALUES (?, ?, ?, datetime('now'))`
		}

		stm, err = DB.Prepare(sqlStm)
		if err != nil {
			return []int{}, err
		}
		defer stm.Close()
		_, err = stm.Exec(userID, id, ltypei)
	} else {
		// User has already liked or disliked the entity
		if existingLikeType == ltypei {
			// If the new like_type is the same as the existing one, remove the like
			if entityType == "post" {
				sqlStm = `DELETE FROM likes WHERE user_id = ? AND post_id = ?`
			} else if entityType == "comment" {
				sqlStm = `DELETE FROM likes WHERE user_id = ? AND comment_id = ?`
			}
		} else {
			// If the new like_type is different, update the like_type
			if entityType == "post" {
				sqlStm = `UPDATE likes SET like_type = ?, created_at = datetime('now') WHERE user_id = ? AND post_id = ?`
			} else if entityType == "comment" {
				sqlStm = `UPDATE likes SET like_type = ?, created_at = datetime('now') WHERE user_id = ? AND comment_id = ?`
			}
		}

		stm, err = DB.Prepare(sqlStm)
		if err != nil {
			return []int{}, err
		}
		defer stm.Close()

		if existingLikeType == ltypei {
			_, err = stm.Exec(userID, id)
		} else {
			_, err = stm.Exec(ltypei, userID, id)
		}
	}

	if err != nil {
		return []int{}, err
	}

	// get the likes and dislikes for this entityType
	likes, err := DBGetLikes(entityType, id)
	if err != nil {
		return []int{}, err
	}

	// update the post or comment likes and dislikes columns
	if entityType == "post" {
		sqlStm = `UPDATE posts SET likes = ?, dislikes = ? WHERE id = ?`
	} else if entityType == "comment" {
		sqlStm = `UPDATE comments SET likes = ?, dislikes = ? WHERE id = ?`
	} else {
		return []int{}, fmt.Errorf("invalid entity type: %s", entityType)
	}

	stm, err = DB.Prepare(sqlStm)
	if err != nil {
		return []int{}, err
	}

	defer stm.Close()

	_, err = stm.Exec(likes[0], likes[1], id)
	if err != nil {
		return []int{}, err
	}

	return likes, nil
}

func DBGetLikes(entityType string, id int) ([]int, error) {

	var likes, dislikes int

	var sqlStm string
	if entityType == "post" {
		sqlStm = `SELECT COUNT(*) FROM likes WHERE post_id = ? AND like_type = 1`
	} else if entityType == "comment" {
		sqlStm = `SELECT COUNT(*) FROM likes WHERE comment_id = ? AND like_type = 1`
	} else {
		return []int{}, fmt.Errorf("invalid entity type: %s", entityType)
	}

	stm, err := DB.Prepare(sqlStm)
	if err != nil {
		return []int{}, err
	}

	defer stm.Close()

	err = stm.QueryRow(id).Scan(&likes)
	if err != nil {
		return []int{}, err
	}

	if entityType == "post" {
		sqlStm = `SELECT COUNT(*) FROM likes WHERE post_id = ? AND like_type = -1`
	} else if entityType == "comment" {
		sqlStm = `SELECT COUNT(*) FROM likes WHERE comment_id = ? AND like_type = -1`
	} else {
		return []int{}, fmt.Errorf("invalid entity type: %s", entityType)
	}

	stm, err = DB.Prepare(sqlStm)
	if err != nil {
		return []int{}, err
	}

	defer stm.Close()

	err = stm.QueryRow(id).Scan(&dislikes)
	if err != nil {
		return []int{}, err
	}

	likesArr := []int{likes, dislikes}
	return likesArr, nil
}

// get likes table data
func GetLikesTable() ([]models.Like, error) {
	sqlStm := `SELECT * FROM likes`
	rows, err := DB.Query(sqlStm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.LikeType, &like.CreatedAt)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

// func RefreshPostData(postID int) (models.Post, error) {
// 	var post models.Post
// 	query := `SELECT id, user_id, title, content, created_at, likes FROM posts WHERE id = ?`
// 	err := DB.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes)
// 	if err != nil {
// 		return post, err
// 	}
// 	return post, nil
// }
