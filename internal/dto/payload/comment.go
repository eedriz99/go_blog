package dto

type CreateCommentPayload struct {
	Content string `json:"content"`
	PostID  string `json:"post_id"`
}

type UpdateCommentPayload struct {
	ID string `json:"id"`
	// PostID  string `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type DeleteCommentPayload struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}
