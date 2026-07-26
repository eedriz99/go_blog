package response

import "github.com/eedriz99/go_blog/internal/model"

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserWithInvitationResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	Username        string `json:"username"`
	InvitationToken string `json:"invitation_token"`
}

func NewUserWithInvitationResponse(user *model.User, invitationToken string) *UserWithInvitationResponse {
	return &UserWithInvitationResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		Username:        user.Username,
		InvitationToken: invitationToken,
	}
}

func NewUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Name:     user.FirstName + " " + user.LastName,
		Username: user.Username,
	}
}
