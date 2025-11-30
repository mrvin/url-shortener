package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	log "github.com/mrvin/url-shortener/internal/logger"
	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockURLsGetter struct {
	mock.Mock
}

func (m *MockURLsGetter) GetURLs(_ context.Context, username string, limit, offset uint64) ([]storage.URL, uint64, error) {
	args := m.Called(username, limit, offset)
	return args.Get(0).([]storage.URL), args.Get(1).(uint64), args.Error(2)
}

func TestGetURLs(t *testing.T) {
	tests := []struct {
		TestName                 string
		Username                 string
		Limit                    uint64
		Offset                   uint64
		StatusCode               int
		URLs                     []storage.URL
		Total                    uint64
		Error                    error
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:   "Success smoke test",
			Username:   "Bob",
			Limit:      10,
			Offset:     1,
			StatusCode: http.StatusOK,
			URLs: []storage.URL{
				{
					URL:       "https://en.wikipedia.org/wiki/Systems_design",
					Alias:     "zn9edcu",
					Count:     24812,
					CreatedAt: time.Date(2025, time.November, 29, 10, 0, 0, 0, time.UTC),
				},
			},
			Total:                    2,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error internal",
			Username:                 "Alice",
			Limit:                    10,
			Offset:                   0,
			StatusCode:               http.StatusInternalServerError,
			URLs:                     nil,
			Total:                    0,
			Error:                    errors.New("internal"),
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "getting urls from storage: internal",
		},
	}
	mockURLsGetter := new(MockURLsGetter)
	handler := NewGetURLs(mockURLsGetter)
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			ctx := log.WithUsername(context.Background(), test.Username)
			url := fmt.Sprintf("/api/urls?limit=%d&offset=%d", test.Limit, test.Offset)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}

			mockURLsGetter.On("GetURLs", test.Username, test.Limit, test.Offset).Return(test.URLs, test.Total, test.Error)

			handler.ServeHTTP(res, req)

			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if res.Code == http.StatusOK {
				var response ResponseGetURLs
				json.Unmarshal(res.Body.Bytes(), &response)
				if !slices.Equal(response.URLs, test.URLs) {
					t.Errorf("expected urls %v but received %v", test.URLs, response.URLs)
				}
				if response.Total != test.Total {
					t.Errorf("expected total %d but received %d", test.Total, response.Total)
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
