package main

import (
	"log"
	"net/http"

	"github.com/blakewilliams/boston-foodie/internal/server"
)

func main() {
	app := server.New()

	logger := log.Default()
	logger.Println("Listening on http://localhost:9090")
	http.ListenAndServe(":9090", app)
}
