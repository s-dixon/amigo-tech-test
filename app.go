package main

import (
	"amigo-tech-test/util"
	"amigo-tech-test/service"
	"os"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"database/sql"
	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func (a *App) Initialise(router service.ServiceRouter, dbConnector util.DatabaseConnector, dbUser, dbPassword, dbName string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	var err error
	a.db, err = dbConnector.Open(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	sr := router.NewServiceRouter(a.db)
	a.router = handlers.LoggingHandler(os.Stdout, sr)
}

func (a *App) Run(httpServer util.HttpServer, addr string) {
	log.Fatal(httpServer.ListenAndServe(addr, a.router))
}