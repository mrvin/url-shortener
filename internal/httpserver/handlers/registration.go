package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *storage.User) error
}

//nolint:tagliatelle
type RequestRegistration struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func NewRegistration(creator UserCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestRegistration

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), "Registration: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), "Registration: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			err := fmt.Errorf("generate hash password: %w", err)
			slog.ErrorContext(req.Context(), "Registration: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		user := storage.User{
			Name:         request.UserName,
			HashPassword: string(hashPassword),
			Role:         "user",
		}

		if err = creator.CreateUser(req.Context(), &user); err != nil {
			err := fmt.Errorf("saving user to storage: %w", err)
			slog.ErrorContext(req.Context(), "Registration: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		httpresponse.WriteOK(res, http.StatusCreated)

		slog.InfoContext(req.Context(), "New user registration was successful")
	}
}
