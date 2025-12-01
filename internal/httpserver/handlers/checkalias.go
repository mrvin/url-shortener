package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrvin/url-shortener/internal/logger"
)

type ResponseCheckAlias struct {
	Exists bool   `json:"exists"`
	Status string `json:"status"`
}

type AliasChecker interface {
	CheckAlias(ctx context.Context, alias string) (bool, error)
}

func NewCheckAlias(checker AliasChecker) HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) (context.Context, int, error) {
		alias := req.PathValue("alias")
		ctx := logger.WithAlias(req.Context(), alias)

		exists, err := checker.CheckAlias(ctx, alias)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("check alias in storage: %w", err)
		}

		// Write json response
		response := ResponseCheckAlias{
			Exists: exists,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("marshal response: %w", err)
		}
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := res.Write(jsonResponse); err != nil {
			return ctx, http.StatusInternalServerError, fmt.Errorf("write response: %w", err)
		}

		return ctx, http.StatusOK, nil
	}
}
