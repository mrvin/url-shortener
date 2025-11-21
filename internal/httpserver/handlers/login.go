package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type UserGetter interface {
	GetUser(ctx context.Context, name string) (*storage.User, error)
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLogin(getter UserGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestLogin

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), "Login: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), "Login: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := getter.GetUser(req.Context(), request.Username)
		if err != nil {
			err := fmt.Errorf("get user: %w", err)
			slog.ErrorContext(req.Context(), "Login: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(request.Password)); err != nil {
			err := fmt.Errorf("compare hash and password: %w", err)
			slog.ErrorContext(req.Context(), "Login: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusUnauthorized)
			return
		}

		// Write json response
		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "User has logged")
	}
}
