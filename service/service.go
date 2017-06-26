package service

import (
	"net/http"
	"database/sql"
	"io/ioutil"
	"github.com/gorilla/mux"
	"strconv"
	"amigo-tech-test/service/model"
)

type ServiceRouter interface {
	NewServiceRouter( db *sql.DB) *mux.Router
}

type MessageServiceRouter struct {
	db *sql.DB
}

func (sr *MessageServiceRouter) NewServiceRouter(db *sql.DB) *mux.Router {
	sr.db = db

	r := mux.NewRouter()
	r.HandleFunc("/messages/", sr.getMessages).Methods("GET")
	r.HandleFunc("/messages/", sr.createMessage).Methods("POST")
	r.HandleFunc("/messages/{Id:[0-9]+}", sr.getMessage).Methods("GET")
	r.HandleFunc("/messages/{Id:[0-9]+}", sr.deleteMessage).Methods("DELETE")

	return r
}

func (sr *MessageServiceRouter) getMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	m := model.Message{Id: id}
	if err := m.GetMessage(sr.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Message not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithString(w, http.StatusOK, m.Value)
}

func (sr *MessageServiceRouter) getMessages(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()

	limit, _  := strconv.Atoi(getQueryParamOrDefault(queryVals, "limit", "20"))
	offset, _  := strconv.Atoi(getQueryParamOrDefault(queryVals, "offset", "0"))
	messageQuery := getQueryParamOrDefault(queryVals, "message", "")
	ipAddress := getQueryParamOrDefault(queryVals, "ip", "")

	if limit < 1 || limit > 20 {
		limit = 20
	}

	if(offset < 0) {
		offset = 0
	}

	result, err := model.Search(sr.db, offset, limit, messageQuery, ipAddress)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

func (sr *MessageServiceRouter) createMessage(w http.ResponseWriter, r *http.Request) {
	var m model.Message
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message payload")
		return
	}
	defer r.Body.Close()

	m.Value = string(bodyBytes)
	ip, error := getClientIp(r)
	if error == nil {
		m.IpAddress = ip.String()
	}

	if err := m.CreateMessage(sr.db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": m.Id})
}

func (sr *MessageServiceRouter) deleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	p := model.Message{Id: id}
	if err := p.DeleteMessage(sr.db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
