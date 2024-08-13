package functions

import (
	"html/template"
	"net/http"
	"strconv"
	"forumProject/internal/database"
	"forumProject/internal/models"
	"forumProject/internal/handlers"
)

func FilterByCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	posts, err := database.GetPosts(0, "ALL")
	if err != nil {
		http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
		return
	}

	selectedCategories := r.Form["category"]
	categoryIDs := make([]int, len(selectedCategories))
	for i, idStr := range selectedCategories {
		categoryIDs[i], _ = strconv.Atoi(idStr)
	}

	var filteredPosts []models.Post
	for _, post := range posts {
		for _, category := range post.Categories {
			if categoryIDExists(categoryIDs, category.ID) {
				filteredPosts = append(filteredPosts, post)
				break
			}
		}
	}

	categories, err := database.GetCategories()
	if err != nil {
		http.Error(w, "Unable to retrieve categories", http.StatusInternalServerError)
		return
	}

	if len(selectedCategories) == 0 {
		filteredPosts = posts
	}

	data := struct {
		Posts      []models.Post
		Categories []models.Category
		LoggedIn   bool
	}{
		Posts:      filteredPosts,
		Categories: categories,
		LoggedIn:   true, // Adjust this based on your session management logic
	}

	t, err := template.ParseFiles("web/index.html", "web/base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func categoryIDExists(ids []int, targetID int) bool {
	for _, id := range ids {
		if id == targetID {
			return true
		}
	}
	return false
}


func FilterByLikes(w http.ResponseWriter, r *http.Request) {
    userID, err := handlers.SessionActive(r)
    if err != nil {
        http.Error(w, "User not logged in", http.StatusUnauthorized)
        return
    }

    showLiked := r.URL.Query().Get("showLiked") == "true"
    showCreated := r.URL.Query().Get("showCreated") == "true"

    var filteredPosts []models.Post

    if !showLiked && !showCreated {
        // If no filter is selected, show all posts
        filteredPosts, err = database.GetPosts(0, "ALL")
    } else {
        allPosts, err := database.GetPosts(0, "ALL")
        if err != nil {
            http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
            return
        }

        userLikes, err := database.GetUserLikes(userID, 1) // Get liked posts
        if err != nil {
            http.Error(w, "Unable to retrieve user likes", http.StatusInternalServerError)
            return
        }

        for _, post := range allPosts {
            if (showLiked && userLikes[post.ID]) || (showCreated && post.UserID == userID) {
                filteredPosts = append(filteredPosts, post)
            }
        }
    }

    if err != nil {
        http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
        return
    }

    categories, err := database.GetCategories()
    if err != nil {
        http.Error(w, "Unable to retrieve categories", http.StatusInternalServerError)
        return
    }

    data := struct {
        Posts      []models.Post
        Categories []models.Category
        LoggedIn   bool
    }{
        Posts:      filteredPosts,
        Categories: categories,
        LoggedIn:   true,
    }

    t, err := template.ParseFiles("web/index.html", "web/base.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

