package database

import (
	"database/sql"

	"forumProject/internal/models"
)

func GetPosts(postID int, typeOf string) ([]models.Post, error) {
	sqlStm := `SELECT posts.id,	
    posts.user_id, posts.title, 
    posts.content, posts.likes, posts.dislikes,
    posts.created_at, users.username
    FROM posts
    JOIN users ON posts.user_id = users.id`

	var rows *sql.Rows
	var err error

	if typeOf == "SINGLE" {
		sqlStm += " WHERE posts.id = ?"
		rows, err = DB.Query(sqlStm, postID)
	} else {
		sqlStm += " ORDER BY posts.created_at DESC" // Order by latest created
		rows, err = DB.Query(sqlStm)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}

		// Get categories for the post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

func CreatePost(postData models.Post) (int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, strftime('%Y-%m-%d %H:%M:%S', 'now', '+3 hours'))",
		postData.UserID, postData.Title, postData.Content)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, category := range postData.Categories {
		_, err = tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, category.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func GetPostsByUser(userID int) ([]models.Post, error) {
	sqlStm := `SELECT posts.id,	
    posts.user_id, posts.title, 
    posts.content, posts.likes, 
    posts.created_at, users.username
    FROM posts
    JOIN users ON posts.user_id = users.id
    WHERE posts.user_id = ?
	ORDER BY posts.created_at DESC`

	rows, err := DB.Query(sqlStm, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsByCategory(category string) ([]models.Post, error) {
	sqlStm := `SELECT posts.id,
        posts.user_id, posts.title,
        posts.content, posts.likes,
        posts.created_at, users.username
    FROM posts
    JOIN users ON posts.user_id = users.id
    JOIN post_categories ON posts.id = post_categories.post_id
    JOIN categories ON post_categories.category_id = categories.id
    WHERE categories.name = ?
	ORDER BY posts.created_at DESC`

	rows, err := DB.Query(sqlStm, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
