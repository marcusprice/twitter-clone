package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/marcusprice/twitter-clone/internal/api"
	"github.com/marcusprice/twitter-clone/internal/constants"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/util"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	util.LoadEnvVariables()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")
	env := os.Getenv("ENV")
	envs := []string{constants.DEV_ENV, constants.PROD_ENV}

	if env == "" {
		panic("ENV environment variable required")
	} else if !slices.Contains(envs, env) {
		panic("Invalid ENV environment variable set, must be 'DEVELOPMENT' or 'PRODUCTION'")
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == "" {
		port = "42069"
	}

	if dbPath == "" {
		panic("DB path required")
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("could not enable foreign keys:", err)
	}
	handler := api.RegisterHandlers(conn)

	logger.LogInfo(fmt.Sprintf("CORE APP LISTENING AT %s:%s", host, port))
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%s", host, port),
			api.Logger(
				api.WithCORS(handler),
			),
		),
	)
}
