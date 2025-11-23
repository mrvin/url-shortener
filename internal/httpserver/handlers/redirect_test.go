package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockGetter struct {
	mock.Mock
}

func (m *MockGetter) GetURL(_ context.Context, alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		TestName                 string
		Alias                    string
		StatusCode               int
		URL                      string
		Error                    error
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:                 "Success smoke test",
			Alias:                    "yc",
			URL:                      "https://yandex.cloud/ru",
			StatusCode:               http.StatusFound,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error alias not found",
			Alias:                    "zn9edcu",
			URL:                      "https://en.wikipedia.org/wiki/Systems_design",
			StatusCode:               http.StatusNotFound,
			Error:                    storage.ErrAliasNotFound,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "alias not found",
		},
		{
			TestName:                 "Error internal",
			Alias:                    "g",
			URL:                      "https://www.google.com/",
			StatusCode:               http.StatusInternalServerError,
			Error:                    errors.New("internal"),
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "failed get url: internal",
		},
	}
	mockGetter := new(MockGetter)
	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodGet+" /{alias...}", NewRedirect(mockGetter))
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()
			res := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+test.Alias, nil)
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}
			mockGetter.On("GetURL", test.Alias).Return(test.URL, test.Error)
			mux.ServeHTTP(res, req)
			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if test.StatusCode == http.StatusFound {
				body, _ := io.ReadAll(res.Body)
				if !bytes.Contains(body, []byte(test.URL)) {
					t.Errorf(`response does not contain "%s"`, test.URL)
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
