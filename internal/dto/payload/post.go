package dto

type CreatePostPayload struct {
	UerID   string   `json:"user_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type DeletePostPayload struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type UpdatePostPayload struct {
	ID      string    `json:"id"`
	Title   *string   `json:"title,omitempty"`
	Content *string   `json:"content,omitempty"`
	Tags    *[]string `json:"tags,omitempty"`
}
