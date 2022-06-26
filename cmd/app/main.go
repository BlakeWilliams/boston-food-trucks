package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blakewilliams/boston-foodie/internal/server"
)

func main() {
	env := server.EnvDevelopment
	if os.Getenv("MEDIUM_ENV") == server.EnvProd {
		env = server.EnvProd
	}

	app := server.New(env)

	logger := log.Default()
	logger.Println("Listening on http://localhost:9090")
	http.ListenAndServe(":9090", app)
}
