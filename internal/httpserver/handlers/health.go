package handlers

import (
	"net/http"

	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
)

func Health(res http.ResponseWriter, _ *http.Request) {
	httpresponse.WriteOK(res, http.StatusOK)
}
