package response

import "github.com/eedriz99/go_blog/internal/model"

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Name:     user.FirstName + " " + user.LastName,
		Username: user.Username,
	}
}
