package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
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

func NewSaveURL(creator URLCreator) HandlerFunc {
	validate := validator.New()
	// Base62 and '_', '-'
	myBase64Regex := regexp.MustCompile("^[0-9a-zA-Z_-]+$")
	err := validate.RegisterValidation("mybase64",
		func(fl validator.FieldLevel) bool {
			return myBase64Regex.MatchString(fl.Field().String())
		})
	if err != nil {
		panic(fmt.Errorf("register validation: %w", err))
	}
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		ctx := req.Context()

		// Read json request
		var request RequestSaveURL
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			return ctx, http.StatusBadRequest, fmt.Errorf("read body request: %w", err)
		}
		if err := json.Unmarshal(body, &request); err != nil {
			return ctx, http.StatusBadRequest, fmt.Errorf("unmarshal body request: %w", err)
		}
		ctx = logger.WithAlias(ctx, request.Alias)
		ctx = logger.WithURL(ctx, request.URL)

		// Validation
		if err := validate.Struct(request); err != nil {
			var vErrors validator.ValidationErrors
			if errors.As(err, &vErrors) {
				return ctx, http.StatusBadRequest, fmt.Errorf("invalid request: tag: %s value: %s", vErrors[0].Tag(), vErrors[0].Value())
			}
			return ctx, http.StatusInternalServerError, fmt.Errorf("validation: %w", err)
		}

		username, err := logger.GetUsernameFromCtx(ctx)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("get user name from ctx: %w", err)
		}

		if err := creator.CreateURL(ctx, username, request.URL, request.Alias); err != nil {
			err := fmt.Errorf("saving url to storage: %w", err)
			if errors.Is(err, storage.ErrAliasExists) {
				return ctx, http.StatusConflict, err
			}
			return ctx, http.StatusInternalServerError, err
		}

		// Write json response
		response := ResponseSaveURL{
			Alias:  request.Alias,
			Status: "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("marshal response: %w", err)
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("write response: %w", err)
		}

		return ctx, http.StatusCreated, nil
	}
}
