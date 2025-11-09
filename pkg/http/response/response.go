package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func WriteOK(res http.ResponseWriter, status int) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		err := fmt.Errorf("WriteOK: marshal error: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	if _, err := res.Write(jsonResponse); err != nil {
		err := fmt.Errorf("WriteOK: write error: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func WriteError(res http.ResponseWriter, msgError string, status int) {
	response := struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{
		Status: "Error",
		Error:  msgError,
	}

	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		err := fmt.Errorf("WriteError: marshal error: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	if _, err := res.Write(jsonResponse); err != nil {
		err := fmt.Errorf("WriteError: write error: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
