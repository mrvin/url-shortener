package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type ResponseCheckAlias struct {
	Exists bool   `json:"exists"`
	Status string `json:"status"`
}

type AliasChecker interface {
	CheckAlias(ctx context.Context, alias string) (bool, error)
}

func NewCheckAlias(checker AliasChecker) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")

		exists, err := checker.CheckAlias(req.Context(), alias)
		if err != nil {
			slog.ErrorContext(req.Context(), "Check alias: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}

		// Write json response
		response := ResponseCheckAlias{
			Exists: exists,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Check alias: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Check alias: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
