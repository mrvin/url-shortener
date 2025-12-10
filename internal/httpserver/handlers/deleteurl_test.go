package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockDBURLDeleter struct {
	mock.Mock
}

func (m *MockDBURLDeleter) DeleteURL(_ context.Context, username, alias string) error {
	args := m.Called(username, alias)
	return args.Error(0)
}

type MockCacheURLDeleter struct {
	mock.Mock
}

func (m *MockCacheURLDeleter) DeleteURL(_ context.Context, alias string) error {
	args := m.Called(alias)
	return args.Error(0)
}

func TestDeleteURL(t *testing.T) {
	tests := []struct {
		TestName                 string
		Username                 string
		Alias                    string
		StatusCode               int
		Error                    error
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:                 "Success smoke test",
			Username:                 "Bob",
			Alias:                    "zn9edcu",
			StatusCode:               http.StatusOK,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error alias not found",
			Username:                 "Alice",
			Alias:                    "systems_design",
			StatusCode:               http.StatusNotFound,
			Error:                    storage.ErrAliasNotFound,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "deleting url from storage: alias not found",
		},
		{
			TestName:                 "Error internal",
			Username:                 "Alice",
			Alias:                    "yc",
			StatusCode:               http.StatusInternalServerError,
			Error:                    errors.New("internal"),
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "deleting url from storage: internal",
		},
	}

	mockDBURLDeleter := new(MockDBURLDeleter)
	mockCacheURLDeleter := new(MockCacheURLDeleter)
	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodDelete+" /api/urls/{alias...}", ErrorHandler("Delete url", NewDeleteURL(mockDBURLDeleter, mockCacheURLDeleter)))
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			ctx := log.WithUsername(context.Background(), test.Username)
			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/api/urls/"+test.Alias, nil)
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}

			mockDBURLDeleter.On("DeleteURL", test.Username, test.Alias).Return(test.Error)
			mockCacheURLDeleter.On("DeleteURL", test.Alias).Return(nil)

			mux.ServeHTTP(res, req)

			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if test.StatusCode == http.StatusOK {
				var response httpresponse.RequestOK
				json.Unmarshal(res.Body.Bytes(), &response)
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
