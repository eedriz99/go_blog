package main

import (
	"log"
	"net/http"
)

// @Summary     Health check
// @Description Returns API status, environment and version
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
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
