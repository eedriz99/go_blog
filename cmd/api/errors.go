package main

import (
	"log"
	"net/http"
)

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	log.Println("ERROR: ", err.Error())
	writeError(w, http.StatusInternalServerError, errorMessage("Internal server error", msg))
}

func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	log.Println("ERROR: ", err.Error())
	writeError(w, http.StatusBadRequest, errorMessage("Bad Request", msg))
}

func (app *application) ResourceNotFoundError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	log.Printf("ERROR: %v", err.Error())
	writeError(w, http.StatusNotFound, errorMessage("Resource not found", msg))
}

func (app *application) UnauthorizedError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	log.Printf("ERROR: %v", err.Error())
	writeError(w, http.StatusUnauthorized, errorMessage("Unauthorized", msg))
}

func (app *application) ConflictError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	log.Printf("ERROR: %v", err.Error())
	writeError(w, http.StatusConflict, errorMessage("Conflict", msg))
}

func errorMessage(fallback string, msg []string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return fallback
}
