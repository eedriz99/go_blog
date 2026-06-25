package main

import (
	"log"
	"net/http"
)

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("ERROR: ", err.Error())
	writeError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("ERROR: ", err.Error())
	writeError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) ResourceNotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("ERROR: %v", err.Error())
	writeError(w, http.StatusNotFound, "Resource not found")
}
