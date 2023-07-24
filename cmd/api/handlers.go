package main

import (
	"fmt"
	"net/http"
)

func (app *Config) HelloWorld(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Root endpoint hit"),
		Data:    "Hello World",
	}

	app.WriteJson(w, http.StatusAccepted, payload)
}
