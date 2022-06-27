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

	logger := log.Default()
	app := server.New(env, logger)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":9090"
	}

	logger.Printf("Listening on http://localhost:%s", addr)
	http.ListenAndServe(addr, app)
}
