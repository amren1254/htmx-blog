package entity

type Blog struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ImagePath     string `json:"image_path"`
	AuthorId      int    `json:"author_id"`
	PublishedDate string `json:"published_date,omitempty"`
}

type BlogRequest struct {
	AuthorId int `json:"author_id"`
	PostId   int `json:"post_id"`
}
