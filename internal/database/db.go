package database

import (
	"database/sql"
	"forumProject/internal/models"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// connect to the database
func connect() error {
	database, err := sql.Open("sqlite3", "internal/database/forumDB.db")
	if err != nil {
		return err
	}
	DB = database
	return nil
}

func InitDB() {
	// connect
	err := connect()
	if err != nil {
		log.Fatal(err)
	}

	// Read the SQL file
	sqlBytes, err := os.ReadFile("internal/database/dbCreator.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL statements
	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")

	// check if database is empty
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		// insert fake data
		CreateFakeData()
		log.Println("Fake data created successfully")
	}
}

func CreateFakeData() {
	// Create users
	users := []models.User{
		{Username: "alice", Email: "alice@example.com", Password: "alice"},
		{Username: "bob", Email: "bob@example.com", Password: "bob"},
		{Username: "charlie", Email: "charlie@example.com", Password: "charlie"},
		{Username: "david", Email: "david@example.com", Password: "david"},
		{Username: "eve", Email: "eve@example.com", Password: "eve"},
	}

	for _, user := range users {
		err := CreateUser(user)
		if err != nil {
			log.Printf("Error creating user %s: %v", user.Username, err)
		}
	}

	// Insert categories
	categories := []string{"Technology", "Sports", "Politics", "Entertainment", "Science"}
	for _, category := range categories {
		_, err := DB.Exec("INSERT INTO categories (name) VALUES (?)", category)
		if err != nil {
			log.Printf("Error inserting category %s: %v", category, err)
		}
	}

	// Insert posts and link to categories
	posts := []struct {
		Title    string
		Content  string
		Category string
	}{
		{"First Post", "This is the content of the first post", "Technology"},
		{"Sports News", "Latest sports updates", "Sports"},
		{"Political Debate", "Discussing current political issues", "Politics"},
		{"Movie Review", "Review of the latest blockbuster", "Entertainment"},
		{"Scientific Discovery", "New findings in quantum physics", "Science"},
		{"Tech Trends", "Emerging technologies in 2023", "Technology"},
		{"Olympic Results", "Medal tally and highlights", "Sports"},
		{"Election Update", "Results from recent elections", "Politics"},
		{"Celebrity Gossip", "Latest Hollywood news", "Entertainment"},
		{"Space Exploration", "Updates on Mars missions", "Science"},
	}

	for _, post := range posts {
		userID := rand.Intn(len(users)) + 1
		result, err := DB.Exec("INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, strftime('%Y-%m-%d %H:%M:%S', 'now', '+3 hours'))",
			userID, post.Title, post.Content, time.Now())
		if err != nil {
			log.Printf("Error inserting post %s: %v", post.Title, err)
			continue
		}

		postID, _ := result.LastInsertId()

		var categoryID int
		err = DB.QueryRow("SELECT id FROM categories WHERE name = ?", post.Category).Scan(&categoryID)
		if err != nil {
			log.Printf("Error getting category ID for %s: %v", post.Category, err)
			continue
		}

		_, err = DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			log.Printf("Error linking post to category: %v", err)
		}
	}

	// Insert comments
	comments := []string{
		"Great post!",
		"I disagree with this.",
		"Interesting perspective.",
		"Thanks for sharing!",
		"Can you elaborate more?",
	}

	for i := 1; i <= len(users)*2; i++ {
		postID := rand.Intn(len(posts)) + 1
		userID := rand.Intn(len(users)) + 1
		content := comments[rand.Intn(len(comments))]

		_, err := DB.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, strftime('%Y-%m-%d %H:%M:%S', 'now', '+3 hours'))",
			postID, userID, content, time.Now())
		if err != nil {
			log.Printf("Error inserting comment: %v", err)
		}
	}

	// Insert likes and dislikes
	for i := 1; i <= len(users)*2; i++ {
		userID := rand.Intn(len(users)) + 1
		postID := rand.Intn(len(posts)) + 1
		likeType := []int{1, -1}[rand.Intn(2)]

		_, err := DB.Exec("INSERT INTO likes (user_id, post_id, like_type, created_at) VALUES (?, ?, ?, ?)",
			userID, postID, likeType, time.Now())
		if err != nil {
			log.Printf("Error inserting like: %v", err)
		}
	}
}
