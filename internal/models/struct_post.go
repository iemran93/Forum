package models

import "encoding/json"

type Post struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	CreatedAt  string     `json:"created_at"`
	Username   string     `json:"username"`
	Likes      int        `json:"likes"`
	Dislikes   int        `json:"dislikes"`
	Categories []Category `json:"categories"`
}

func (p *Post) UnmarshalJSON(data []byte) error {
	type Alias Post
	aux := &struct {
		Categories []int `json:"categories"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Categories = make([]Category, len(aux.Categories))
	for i, id := range aux.Categories {
		p.Categories[i] = Category{ID: id}
	}
	return nil
}
