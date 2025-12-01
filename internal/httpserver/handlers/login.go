package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrvin/url-shortener/internal/logger"
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

func NewLogin(getter UserGetter) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		ctx := req.Context()

		// Read json request
		var request RequestLogin
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			return ctx, http.StatusBadRequest, fmt.Errorf("read body request: %w", err)
		}
		if err := json.Unmarshal(body, &request); err != nil {
			return ctx, http.StatusBadRequest, fmt.Errorf("unmarshal body request: %w", err)
		}
		ctx = logger.WithUsername(ctx, request.Username)

		user, err := getter.GetUser(ctx, request.Username)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("getting user from storage: %w", err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(request.Password)); err != nil {
			return ctx, http.StatusUnauthorized, fmt.Errorf("compare hash and password: %w", err)
		}

		// Write json response
		httpresponse.WriteOK(res, http.StatusOK)

		return ctx, http.StatusOK, nil
	}
}
