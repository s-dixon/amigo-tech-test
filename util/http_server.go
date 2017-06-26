package util

import (
	"net/http"
)

type HttpServer interface {
	ListenAndServe(addr string, router http.Handler) error
}

type DefaultHttpServer struct {}

func (DefaultHttpServer) ListenAndServe(addr string, router http.Handler) error {
	return http.ListenAndServe(addr, router)
}

