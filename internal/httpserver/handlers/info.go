package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

//nolint:gochecknoglobals
var (
	tag  string
	hash string
	date string
)

type ResponseInfo struct {
	Tag  string `json:"tag"`
	Hash string `json:"hash"`
	Date string `json:"date"`
}

func Info(res http.ResponseWriter, _ *http.Request) {
	response := ResponseInfo{
		Tag:  tag,
		Hash: hash,
		Date: date,
	}
	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		err := fmt.Errorf("Info: marshal response: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := res.Write(jsonResponse); err != nil {
		err := fmt.Errorf("Info: write response: %w", err)
		slog.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
