package main

import (
	"log"
	"net/http"
)

func (app *application) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Internal server error: ", err.Error())
	writeError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Bad Request Error: ", err.Error())
	writeError(w, http.StatusBadRequest, "Error Bad Request")
}

// func (app *application) ResourceNotFoundError(w http.ResponseWriter, r *http.Request, err error) {
// 	// log.Printf()
// 	// writeError(w, http.Stat)
// }
