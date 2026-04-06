package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	payload := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJson(w, http.StatusOK, payload); err != nil {
		err = writeError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Println("Error writing response:", err.Error())
		}
		return
	}
	return
}
