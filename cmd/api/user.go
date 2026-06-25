package main

import (
	"net/http"

	"github.com/eedriz99/go_blog/internal/dto/response"
)

type userKey string

const userContextKey userKey = "user"

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	writeJson(w, http.StatusOK, response.NewUserResponse(user))
}
