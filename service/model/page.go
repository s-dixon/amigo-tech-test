package model

type Page struct {
	Offset int `json:"offset"`
	Limit int  `json:"limit"`
	TotalCount int  `json:"total_count"`
	Results interface{}  `json:"results"`
}
