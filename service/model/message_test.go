package model

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"time"
)

func Test_ShouldRetrieveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	expectedDateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	columns := []string{"value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT value, ip_address, date_created FROM messages").
		WithArgs(123).
		WillReturnRows(sqlmock.NewRows(columns).AddRow("Test message value", "192.168.200.201", expectedDateCreated))

	message := Message{Id: 123}

	err = message.GetMessage(db)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t,"Test message value", message.Value,"Message value was not mapped as expected")
	assert.Equal(t,"192.168.200.201", message.IpAddress, "Message IP address was not mapped as expected")
	assert.Equal(t,"2017-06-25 14:22:12.296925 +0000 UTC", message.DateCreated.String(), "Message date created was not mapped as expected")
}

func Test_ShouldCreateNewMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	columns := []string{"id"}
	mock.ExpectQuery("INSERT INTO messages").
		WithArgs("Test message value", "192.168.200.201").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(22))

	message := Message{Value: "Test message value", IpAddress: "192.168.200.201"}

	err = message.CreateMessage(db)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, 22, message.Id, "Message Id was not mapped as expected")
}

func Test_ShouldDeleteMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mock.ExpectExec("DELETE FROM messages").WithArgs(123).WillReturnResult(sqlmock.NewResult(0,1))
	message := Message{Id: 123}

	err = message.DeleteMessage(db)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func Test_ShouldSearchForMessagesWithoutCriteriaAndMapResultsCorrectly(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	expectedDateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	columns := []string{"id", "value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "Test message value 1", "192.168.200.201", expectedDateCreated).
		AddRow(2, "Test message value 2", "127.0.0.1", expectedDateCreated))

	page, error := Search(db, 0, 10, "", "")

	assert.Nil(t, error)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, 0, page.Offset, "Page offset was not set correctly")
	assert.Equal(t, 10, page.Limit, "Page limit was not set correctly")
	assert.Equal(t, 100, page.TotalCount, "Total count was not set correctly")

	messages := page.Results.([]Message)
	assert.Equal(t, 2, len(messages), "Expected two messages to be returned")

	message1 := messages[0]
	assert.Equal(t, 1, message1.Id, "First message Id was not as expected")
	assert.Equal(t, "Test message value 1", message1.Value, "First message Value was not as expected")
	assert.Equal(t, "192.168.200.201", message1.IpAddress, "First message IpAddress was not as expected")

	message2 := messages[1]
	assert.Equal(t, 2, message2.Id, "Second message Id was not as expected")
	assert.Equal(t, "Test message value 2", message2.Value, "Second message Value was not as expected")
	assert.Equal(t, "127.0.0.1", message2.IpAddress, "Second message IpAddress was not as expected")
}

func Test_ShouldSearchForMessagesWithMessageValueCriteria(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	expectedDateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("%Test message value%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	columns := []string{"id", "value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WithArgs("%Test message value%").
		WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "Test message value", "192.168.200.201", expectedDateCreated))

	page, error := Search(db, 0, 10, "Test message value", "")

	assert.Nil(t, error)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, 0, page.Offset, "Page offset was not set correctly")
	assert.Equal(t, 10, page.Limit, "Page limit was not set correctly")
	assert.Equal(t, 100, page.TotalCount, "Total count was not set correctly")
	assert.Equal(t, 1, len(page.Results.([]Message)), "Expected one message to be returned")
}

func Test_ShouldSearchForMessagesWithIpAddressCriteria(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	expectedDateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("192.168%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	columns := []string{"id", "value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WithArgs("192.168%").
		WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "Test message value", "192.168.200.201", expectedDateCreated))

	page, error := Search(db, 0, 10, "", "192.168")

	assert.Nil(t, error)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, 0, page.Offset, "Page offset was not set correctly")
	assert.Equal(t, 10, page.Limit, "Page limit was not set correctly")
	assert.Equal(t, 100, page.TotalCount, "Total count was not set correctly")
	assert.Equal(t, 1, len(page.Results.([]Message)), "Expected one message to be returned")
}

func Test_ShouldSearchForMessagesWithMessageValueAndIpAddressCriteria(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	expectedDateCreated, err := time.Parse(time.RFC3339, "2017-06-25T14:22:12.296925Z")

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("%Test message value%", "192.168%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
		AddRow(100))

	columns := []string{"id", "value", "ip_address", "date_created"}
	mock.ExpectQuery("SELECT id, value, ip_address, date_created FROM messages").
		WithArgs("%Test message value%", "192.168%").
		WillReturnRows(sqlmock.NewRows(columns).
		AddRow(1, "Test message value", "192.168.200.201", expectedDateCreated))

	page, error := Search(db, 0, 10, "Test message value", "192.168")

	assert.Nil(t, error)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

	assert.Equal(t, 0, page.Offset, "Page offset was not set correctly")
	assert.Equal(t, 10, page.Limit, "Page limit was not set correctly")
	assert.Equal(t, 100, page.TotalCount, "Total count was not set correctly")
	assert.Equal(t, 1, len(page.Results.([]Message)), "Expected one message to be returned")
}