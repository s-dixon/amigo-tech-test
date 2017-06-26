package service

import (
	"net/http"
	"encoding/json"
	"net"
	"fmt"
	"net/url"
)

func getQueryParamOrDefault(queryVals url.Values, key, defaultVal string) string{
	if v := queryVals[key]; len(v) == 1 {
		return v[0]
	}
	return defaultVal
}

func getClientIp(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("Client IP: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("Client IP: %q is not a valid IP address", ip)
	}
	return userIP, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithString(w http.ResponseWriter, code int, text string) {
	w.WriteHeader(code)
	w.Write([]byte(text))
}