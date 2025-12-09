package handlers

import (
	"net/http"
)

func NewAPIDocs(path string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(res, req, path)
	}
}
