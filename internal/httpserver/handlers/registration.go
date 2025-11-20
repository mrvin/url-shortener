package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *storage.User) error
}

type RequestRegistration struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func NewRegistration(creator UserCreator) http.HandlerFunc {
	validate := validator.New()
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

		// Validation
		if err := validate.Struct(request); err != nil {
			var vErrors validator.ValidationErrors
			if errors.As(err, &vErrors) {
				err := fmt.Errorf("invalid request: tag: %s value: %s", vErrors[0].Tag(), vErrors[0].Value())
				slog.ErrorContext(req.Context(), "Save url: "+err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
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
			Name:         request.Username,
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
