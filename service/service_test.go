package service

import (
	"testing"
	"net/http/httptest"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"errors"
	"bytes"
	"time"
)

var router *mux.Router

func Test_ShouldGetMessageWithoutErrors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)
	dateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	columns := []string{"value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT value, ip_address, date_created FROM messages").
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows(columns).AddRow("Test message value", "192.168.200.201", dateCreated))

	req, _ := http.NewRequest("GET", "/messages/11", nil)
	response := executeRequest(req)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "Test message value",response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToGetMessageDueToMessageNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectQuery("SELECT value, ip_address, date_created FROM messages").
		WithArgs(11).
		WillReturnError(sql.ErrNoRows)

	req, _ := http.NewRequest("GET", "/messages/11", nil)
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"error\":\"Message not found\"}",response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToGetMessageDueToError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectQuery("SELECT value, ip_address, date_created FROM messages").
		WithArgs(11).
		WillReturnError(errors.New("Database connection closed"))

	req, _ := http.NewRequest("GET", "/messages/11", nil)
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"error\":\"Database connection closed\"}",response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToGetMessageDueToNumberError(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	req, _ := http.NewRequest("GET", "/messages/999999999999999999999999999999999999999999", nil)
	response := executeRequest(req)

	assert.Equal(t, "{\"error\":\"Invalid message ID\"}",response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldSuccessfullyRetrieveMessages(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)
	dateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("%Test message value%", "192.168%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	columns := []string{"id", "value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WithArgs("%Test message value%", "192.168%").
		WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "Test message value", "192.168.200.201", dateCreated))

	req, _ := http.NewRequest("GET", "/messages/?message=Test%20message%20value&ip=192.168", nil)
	response := executeRequest(req)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"offset\":0,\"limit\":20,\"total_count\":100,\"results\":[{\"id\":1,\"value\":\"Test message value\",\"ip_address\":\"192.168.200.201\",\"date_created\":\"2017-06-25T14:22:12.296925Z\"}]}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToRetrieveMessagesAndRespondWithError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("%Test message value%", "192.168%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WithArgs("%Test message value%", "192.168%").
		WillReturnError(errors.New("Database connection closed"))

	req, _ := http.NewRequest("GET", "/messages/?message=Test%20message%20value&ip=192.168", nil)
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"error\":\"Database connection closed\"}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldSuccessfullyCreateMessageWithClientIpAddress(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	columns := []string{"id"}
	mock.ExpectQuery("INSERT INTO messages").
		WithArgs("Test message value", "192.168.200.201").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(22))

	req, _ := http.NewRequest("POST", "/messages/", bytes.NewBuffer([]byte("Test message value")))
	req.RemoteAddr = "192.168.200.201:5000"
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"id\":22}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldSuccessfullyCreateMessageWithoutClientIpAddress(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	columns := []string{"id"}
	mock.ExpectQuery("INSERT INTO messages").
		WithArgs("Test message value", "").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(22))

	req, _ := http.NewRequest("POST", "/messages/", bytes.NewBuffer([]byte("Test message value")))
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, "{\"id\":22}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToCreateMessageDueToError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectQuery("INSERT INTO messages").
		WithArgs("Test message value", "").
		WillReturnError(errors.New("Database connection closed"))

	req, _ := http.NewRequest("POST", "/messages/", bytes.NewBuffer([]byte("Test message value")))
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t,"{\"error\":\"Database connection closed\"}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldDeleteMessageSuccessfully(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectExec("DELETE FROM messages").
		WithArgs(123).
		WillReturnResult(sqlmock.NewResult(0,1))

	req, _ := http.NewRequest("DELETE", "/messages/123", nil)
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t,"{\"result\":\"success\"}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToDeleteMessageDueToInvalidMessageId(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	req, _ := http.NewRequest("DELETE", "/messages/999999999999999999999999999999999999999999", nil)
	response := executeRequest(req)

	assert.Equal(t,"{\"error\":\"Invalid message ID\"}", response.Body.String(), "Response body does not match expected value")
}

func Test_ShouldFailToDeleteMessageDueToErrorDeletingMessageFromDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	router = (&MessageServiceRouter{}).NewServiceRouter(db)

	mock.ExpectExec("DELETE FROM messages").
		WithArgs(123).
		WillReturnError(errors.New("Database connection closed"))

	req, _ := http.NewRequest("DELETE", "/messages/123", nil)
	response := executeRequest(req)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t,"{\"error\":\"Database connection closed\"}", response.Body.String(), "Response body does not match expected value")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}