package main

import (
	"net/http"

	"github.com/eedriz99/go_blog/internal/dto/response"
)

type userKey string

const userContextKey userKey = "user"

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

}

// @Summary     Read user
// @Description Read a user by ID
// @Tags        users
// @Produce     json
// @Param       userID path string true "User ID"
// @Success     200 {object} response.UserResponse
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /users/{userID} [get]
func (app *application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	writeJson(w, http.StatusOK, response.NewUserResponse(user))
}
