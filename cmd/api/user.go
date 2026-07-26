package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/dto/response"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type userKey string

const userContextKey userKey = "user"

// @Summary     Register user
// @Description Create a new user account and send an activation invitation
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       payload body payload.CreateUserPayload true "User registration payload"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/register [post]
func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// create user and hash password
	var payload payload.CreateUserPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	user := &model.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Password:  string(hash),
	}

	token, err := app.store.Users.CreateWithInvitation(ctx, user, app.config.mail.expiry)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case "users_email_key":
				app.ConflictError(w, r, err, "User with this email already exists")
				return
			case "users_username_key":
				app.ConflictError(w, r, err, "User with this username already exists")
				return
			default:
				app.InternalServerError(w, r, err, "Failed to create user")
				return
			}
		} else {
			app.InternalServerError(w, r, err)
			return
		}

	}

	// TODO: Email the activation token to the user registered Email

	log.Printf("Token: %v\n", &token)
	writeJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
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
	writeJSON(w, http.StatusOK, response.NewUserResponse(user))
}

// @Summary     Login
// @Description Authenticate a user with email and password
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       payload body payload.LoginPayload true "Login payload"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/login [post]
func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get user credentials from request body
	var payload payload.LoginPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch user from database by email
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, store.ErrorNotFound) {
			app.UnauthorizedError(w, r, err, "Invalid email or password")
		} else {
			app.InternalServerError(w, r, err, "Failed to fetch user")
		}
		return
	}

	// Validate credentials
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		app.UnauthorizedError(w, r, err, "Invalid email or password")
		return
	} else if !user.IsActive {
		app.UnauthorizedError(w, r, errors.New("account not activated"), "Not an active user")
		return
	}

	// TODO: Reroute to home page if valid, otherwise return an error response
	writeJSON(w, http.StatusOK, "login successful")
}

// @Summary     Activate user
// @Description Activate a user account using its invitation token
// @Tags        auth
// @Produce     json
// @Param       token path string true "Invitation token"
// @Success     202 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/activate/{token} [post]
func (app *application) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	if err := app.store.Users.ActivateByInvitationToken(ctx, token); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusAccepted, "Activated")

}
