package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
)

type URLCreator interface {
	CreateURL(ctx context.Context, username, url, alias string) error
}

type RequestSaveURL struct {
	URL   string `json:"url"   validate:"required,url"`
	Alias string `json:"alias" validate:"required,mybase64"`
}

type ResponseSaveURL struct {
	Alias  string `json:"alias"`
	Status string `json:"status"`
}

func NewSaveURL(creator URLCreator) http.HandlerFunc {
	validate := validator.New()
	// Base62 and '_', '-'
	myBase64Regex := regexp.MustCompile("^[0-9a-zA-Z_-]+$")
	err := validate.RegisterValidation("mybase64",
		func(fl validator.FieldLevel) bool {
			return myBase64Regex.MatchString(fl.Field().String())
		})
	if err != nil {
		panic("Register validation: " + err.Error())
	}
	return func(res http.ResponseWriter, req *http.Request) {
		var request RequestSaveURL

		// Read json request
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
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

		username, err := logger.GetUsernameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := creator.CreateURL(req.Context(), username, request.URL, request.Alias); err != nil {
			if errors.Is(err, storage.ErrAliasExists) {
				slog.ErrorContext(req.Context(), "Save url: "+err.Error(), slog.String("alias", request.Alias))
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
			err := fmt.Errorf("failed save url: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error(), slog.String("alias", request.Alias))
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseSaveURL{
			Alias:  request.Alias,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(req.Context(), "Create new url",
			slog.String("alias", request.Alias),
			slog.String("url", request.URL),
		)
	}
}
