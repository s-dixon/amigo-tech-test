package model

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"fmt"
	"time"
)

type Message struct {
	Id    int	 `json:"id"`
	Value string `json:"value"`
	IpAddress string `json:"ip_address"`
	DateCreated time.Time `json:"date_created""`
}

func (m *Message) GetMessage(db *sql.DB) error {
	return sq.Select("value", "ip_address", "date_created").
		From("messages").
		Where(sq.Eq{"id": &m.Id}).
		RunWith(db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&m.Value, &m.IpAddress, &m.DateCreated)
}

func (m *Message) DeleteMessage(db *sql.DB) error {
	_, err:= sq.Delete("messages").
		Where(sq.Eq{"id": &m.Id}).
		RunWith(db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	return err
}

func (m *Message) CreateMessage(db *sql.DB) error {
	return sq.Insert("messages").
		Columns("value", "ip_address").
		Values(m.Value, m.IpAddress).
		Suffix("RETURNING \"id\"").
		RunWith(db).
		PlaceholderFormat(sq.Dollar).
		QueryRow().
		Scan(&m.Id)
}

func Search(db *sql.DB, offset, limit int, message, ipAddress string) (*Page, error) {
	pageQuery := sq.Select("id", "value", "ip_address", "date_created").From("messages").RunWith(db)
	countQuery := sq.Select("COUNT(id)").From("messages").RunWith(db)

	if message != "" {
		addWhereCondition("value LIKE ?", fmt.Sprint("%", message, "%"), &pageQuery, &countQuery)
	}

	if ipAddress != "" {
		addWhereCondition("TEXT(ip_address) LIKE ?", fmt.Sprint(ipAddress, "%"), &pageQuery, &countQuery)
	}

	var totalCount int
	countQuery.PlaceholderFormat(sq.Dollar).QueryRow().Scan(&totalCount)

	rows, err :=  pageQuery.
		PlaceholderFormat(sq.Dollar).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Id, &m.Value, &m.IpAddress, &m.DateCreated); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return &Page{offset, limit, totalCount, messages}, nil
}

func addWhereCondition(predicate, value string, queries ...*sq.SelectBuilder) {
	for _, query := range queries {
		where := query.Where(predicate, value)
		*query = where
	}
}
