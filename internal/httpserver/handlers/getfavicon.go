package handlers

import (
	"net/http"
)

func GetFavicon(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/static/favicon.ico", http.StatusMovedPermanently)
}
