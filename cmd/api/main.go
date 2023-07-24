package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "3000"

type Config struct{}

func main() {
	log.Println("Starting Rest Service")

	app := Config{}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
