package handlers

import (
	"fmt"
	"forumProject/internal/database"
	"forumProject/internal/models"
	"net/http"
	"text/template"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var err error

		filterParams := r.URL.Query()
		var filteredPosts []models.Post

		// filter by category
		if _, exist := filterParams["categories"]; exist {
			filteredPosts, err = filterByCategory(filterParams["categories"])
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
				return
			}
		} else {
			filteredPosts, err = database.GetPosts(0, "ALL")
			if err != nil {
				RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
				return
			}
		}

		// filter by user likes/created
		if _, exist := filterParams["byUser"]; exist {
			userID, err := SessionActive(r)
			if err != nil {
				RenderErrorPage(w, http.StatusForbidden, "Forbidden: No session token")
				return
			}

			var userPosts []models.Post
			var userLikes []models.Post
			crposts := false
			likeposts := false

			for _, option := range filterParams["byUser"] {
				switch option {
				case "crposts":
					userPosts, err = database.GetPostsByUser(userID)
					if err != nil {
						RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
						return
					}
					crposts = true
				case "likeposts":
					userLikes, err = filterByUserLiked(userID)
					if err != nil {
						RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
						return
					}
					likeposts = true
				}
			}

			if crposts && likeposts {
				filteredPosts = mergePosts(filteredPosts, userPosts, userLikes)
			} else if crposts {
				filteredPosts = mergePosts(filteredPosts, userPosts, nil)
			} else if likeposts {
				filteredPosts = mergePosts(filteredPosts, nil, userLikes)
			}
		}

		categories, err := database.GetCategories() // Implement this function
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}

		loggedIn := false
		_, err = SessionActive(r)
		if err == nil {
			loggedIn = true
		}

		// render the filtered posts
		pd := PageData{
			Posts:      filteredPosts,
			LoggedIn:   loggedIn,
			Categories: categories,
		}
		// serve the template with the data
		t, err := template.ParseFiles("web/index.html", "web/base.html")
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
		err = t.Execute(w, pd)
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
			return
		}
	}
}

func filterByCategory(categories []string) ([]models.Post, error) {
	var allPosts []models.Post
	for _, category := range categories {
		posts, err := database.GetPostsByCategory(category)
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, posts...)
	}
	return allPosts, nil
}

func filterByUserLiked(userID int) ([]models.Post, error) {
	userLikes := []models.Post{}

	likesData, err := database.GetLikesTable()
	if err != nil {
		return nil, err
	}

	for _, like := range likesData {
		if like.UserID == userID && like.LikeType == 1 {
			post, err := database.GetPosts(*like.PostID, "SINGLE")
			if err != nil {
				return nil, err
			}
			userLikes = append(userLikes, post[0])
		}
	}

	return userLikes, nil
}

func mergePosts(existing, createdPosts, likedPosts []models.Post) []models.Post {
	if createdPosts == nil && likedPosts == nil {
		return existing
	}

	filteredPosts := []models.Post{}

	filterPosts := func(posts []models.Post) {
		for _, post := range posts {
			for _, existingPost := range existing {
				// if the post is already in the existing posts, add it to the filtered posts
				// but only if it's not already in the filtered posts
				if post.ID == existingPost.ID {
					exists := false
					for _, filteredPost := range filteredPosts {
						if filteredPost.ID == existingPost.ID {
							exists = true
							break
						}
					}
					if !exists {
						filteredPosts = append(filteredPosts, existingPost)
					}
				}
			}
		}
	}

	if createdPosts != nil {
		filterPosts(createdPosts)
	}

	if likedPosts != nil {
		filterPosts(likedPosts)
	}

	return filteredPosts
}
