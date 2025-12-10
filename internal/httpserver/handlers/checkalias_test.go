package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockAliasChecker struct {
	mock.Mock
}

func (m *MockAliasChecker) CheckAlias(ctx context.Context, alias string) (bool, error) {
	args := m.Called(alias)
	return args.Bool(0), args.Error(1)
}

func TestCheckAlias(t *testing.T) {
	tests := []struct {
		TestName                 string
		Alias                    string
		StatusCode               int
		Error                    error
		Exists                   bool
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:                 "Alias does not exist",
			Alias:                    "zn9edcu",
			StatusCode:               http.StatusOK,
			Error:                    nil,
			Exists:                   false,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Alias exists",
			Alias:                    "systems_design",
			StatusCode:               http.StatusOK,
			Error:                    nil,
			Exists:                   true,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error internal",
			Alias:                    "yc",
			StatusCode:               http.StatusInternalServerError,
			Error:                    errors.New("internal"),
			Exists:                   false,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "check alias in storage: internal",
		},
	}

	mockAliasChecker := new(MockAliasChecker)
	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodGet+" /api/urls/check/{alias...}", ErrorHandler("Check alias", NewCheckAlias(mockAliasChecker)))
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/urls/check/"+test.Alias, nil)
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}

			mockAliasChecker.On("CheckAlias", test.Alias).Return(test.Exists, test.Error)

			mux.ServeHTTP(res, req)

			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if test.StatusCode == http.StatusOK {
				var response ResponseCheckAlias
				json.Unmarshal(res.Body.Bytes(), &response)
				if response.Status != test.ExpectedStatus {
					t.Errorf(`expected status "%s" but received "%s"`, test.ExpectedStatus, response.Status)
				}
				if response.Exists != test.Exists {
					t.Errorf(`expected exists "%t" but received "%t"`, test.Exists, response.Exists)
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
