package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	// log.Printf("Response Read: %v", data)
	return nil
}

func writeError(w http.ResponseWriter, status int, err string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJson(w, status, &envelope{Error: err})
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	// log.Printf("Payload Read: %v", data)
	return nil
}
