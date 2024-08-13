package models

type Like struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	PostID    *int   `json:"post_id,omitempty"`    // Nullable
	CommentID *int   `json:"comment_id,omitempty"` // Nullable
	LikeType  int    `json:"like_type"`            // 1 for like, -1 for dislike
	CreatedAt string `json:"created_at"`
}
