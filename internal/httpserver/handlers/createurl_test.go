package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockURLCreator struct {
	mock.Mock
}

func (m *MockURLCreator) CreateURL(_ context.Context, username, url, alias string) error {
	args := m.Called(username, url, alias)
	return args.Error(0)
}

func TestCreateURL(t *testing.T) {
	tests := []struct {
		TestName                 string
		Username                 string
		URL                      string
		Alias                    string
		StatusCode               int
		Error                    error
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:                 "Success smoke test",
			Username:                 "Bob",
			URL:                      "https://yandex.cloud/ru",
			Alias:                    "yc",
			StatusCode:               http.StatusCreated,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Success alias with underscore",
			Username:                 "Bob",
			URL:                      "https://en.wikipedia.org/wiki/Systems_design",
			Alias:                    "sys_dsgn",
			StatusCode:               http.StatusCreated,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Success alias with dash",
			Username:                 "Bob",
			URL:                      "https://www.youtube.com/",
			Alias:                    "y-t",
			StatusCode:               http.StatusCreated,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error alias already exists",
			Username:                 "Bob",
			URL:                      "https://www.google.com/",
			Alias:                    "g",
			StatusCode:               http.StatusConflict,
			Error:                    storage.ErrAliasExists,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "saving url to storage: alias already exists",
		},
		{
			TestName:                 "Error invalid url",
			Username:                 "Bob",
			URL:                      "//www.google.com/",
			Alias:                    "g",
			StatusCode:               http.StatusBadRequest,
			Error:                    nil,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "invalid request: tag: url value: //www.google.com/",
		},
		{
			TestName:                 "Error invalid alias",
			Username:                 "Bob",
			URL:                      "https://www.google.com/",
			Alias:                    "api/",
			StatusCode:               http.StatusBadRequest,
			Error:                    nil,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "invalid request: tag: mybase64 value: api/",
		},
	}

	mockCreator := new(MockURLCreator)
	handler := NewSaveURL(mockCreator)
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()
			res := httptest.NewRecorder()
			dataRequest, err := json.Marshal(RequestSaveURL{URL: test.URL, Alias: test.Alias})
			if err != nil {
				t.Fatalf("cant marshal json: %v", err)
			}
			ctx := log.WithUsername(context.Background(), test.Username)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/data/shorten", bytes.NewReader(dataRequest))
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}
			mockCreator.On("CreateURL", test.Username, test.URL, test.Alias).Return(test.Error)
			handler.ServeHTTP(res, req)
			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if test.StatusCode == http.StatusCreated {
				var response ResponseSaveURL
				json.Unmarshal(res.Body.Bytes(), &response)
				if response.Alias != test.Alias {
					t.Errorf(`expected alias "%s" but received "%s"`, test.Alias, response.Alias)
				}
				if response.Status != test.ExpectedStatus {
					t.Errorf(`expected status "%s" but received "%s"`, test.ExpectedStatus, response.Status)
				}
			} else {
				var response httpresponse.RequestError
				json.Unmarshal(res.Body.Bytes(), &response)
				if response.Status != test.ExpectedStatus {
					t.Errorf(`expected status "%s" but received "%s"`, test.ExpectedStatus, response.Status)
				}
				if response.Error != test.ExpectedErrorDescription {
					t.Errorf(`expected description "%s" but received "%s"`, test.ExpectedErrorDescription, response.Error)
				}
			}
		})
	}
}
