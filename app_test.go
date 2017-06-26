package main

import (
	"testing"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
)

type stubDatabaseConnector struct{
	DB *sql.DB
	ConnString string
}
func (dc *stubDatabaseConnector) Open(connectionString string) (*sql.DB, error) {
	dc.ConnString = connectionString
	dc.DB, _, _ = sqlmock.New()
	return dc.DB, nil
}

type stubServiceRouter struct {
	DB *sql.DB
}

func (sr *stubServiceRouter) NewServiceRouter(db *sql.DB) *mux.Router {
	sr.DB = db
	return mux.NewRouter()
}

type stubHttpServer struct {
	Addr string
	Router http.Handler
}

func (hs *stubHttpServer) ListenAndServe(addr string, router http.Handler) error {
	hs.Addr = addr
	hs.Router = router
	return nil
}

func TestApp_Initialisation(t *testing.T){
	app := App{}
	router := &stubServiceRouter{}
	connector := &stubDatabaseConnector{}

	app.Initialise(router, connector, "username", "password", "database")

	assert.Equal(t, "user=username password=password dbname=database sslmode=disable", connector.ConnString, "Connection string does not match the expected value")
	assert.Equal(t, connector.DB, router.DB, "Database object was not injected into the router")
}

func TestApp_Run(t *testing.T){
	app := App{}
	router := &stubServiceRouter{}
	connector := &stubDatabaseConnector{}
	server := &stubHttpServer{}

	app.Initialise(router, connector, "username", "password", "database")
	app.Run(server, ":8080")

	assert.Equal(t, ":8080", server.Addr, "Server address was not set correctly")
	assert.Equal(t, server.Router, app.router, "Http router was not configured correctly")
}