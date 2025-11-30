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

type MockDBURLGetter struct {
	mock.Mock
}

func (m *MockDBURLGetter) GetURL(_ context.Context, alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func (m *MockDBURLGetter) CountIncrement(alias string) error {
	args := m.Called(alias)
	return args.Error(0)
}

type MockCacheURLGetter struct {
	mock.Mock
}

func (m *MockCacheURLGetter) GetURL(_ context.Context, alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func (m *MockCacheURLGetter) SetURL(_ context.Context, url, alias string) error {
	args := m.Called(alias, url)
	return args.Error(0)
}

func TestRedirect(t *testing.T) {
	mockDBURLGetter := new(MockDBURLGetter)
	mockCacheURLGetter := new(MockCacheURLGetter)
	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodGet+" /{alias...}", NewRedirect(mockDBURLGetter, mockCacheURLGetter))

	t.Run("Success smoke test and cache miss", func(t *testing.T) {
		t.Parallel()

		alias := "yc"
		url := "https://yandex.cloud/ru"

		res := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+alias, nil)
		if err != nil {
			t.Fatalf("create new request: %v", err)
		}

		mockCacheURLGetter.On("GetURL", alias).Return("", nil)
		mockDBURLGetter.On("GetURL", alias).Return(url, nil)
		mockCacheURLGetter.On("SetURL", alias, url).Return(nil)

		mux.ServeHTTP(res, req)

		status := http.StatusFound
		if res.Code != status {
			t.Errorf("expected status code %d but received %d", status, res.Code)
		}
		body, _ := io.ReadAll(res.Body)
		if !bytes.Contains(body, []byte(url)) {
			t.Errorf(`response does not contain "%s"`, url)
		}

	})

	t.Run("Success smoke test and cache hit", func(t *testing.T) {
		t.Parallel()

		alias := "zn9edcu"
		url := "https://en.wikipedia.org/wiki/Systems_design"

		res := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+alias, nil)
		if err != nil {
			t.Fatalf("create new request: %v", err)
		}

		mockCacheURLGetter.On("GetURL", alias).Return(url, nil)
		mockDBURLGetter.On("CountIncrement", alias).Return(nil)

		mux.ServeHTTP(res, req)

		status := http.StatusFound
		if res.Code != status {
			t.Errorf("expected status code %d but received %d", status, res.Code)
		}
		body, _ := io.ReadAll(res.Body)
		if !bytes.Contains(body, []byte(url)) {
			t.Errorf(`response does not contain "%s"`, url)
		}
	})

	t.Run("Error alias not found", func(t *testing.T) {
		t.Parallel()

		alias := "7OeLY0"

		res := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+alias, nil)
		if err != nil {
			t.Fatalf("create new request: %v", err)
		}

		mockCacheURLGetter.On("GetURL", alias).Return("", nil)
		mockDBURLGetter.On("GetURL", alias).Return("", storage.ErrAliasNotFound)

		mux.ServeHTTP(res, req)

		status := http.StatusNotFound
		if res.Code != status {
			t.Errorf("expected status code %d but received %d", status, res.Code)
		}

		expectedStatus := "Error"
		expectedErrorDescription := "getting url from storage: alias not found"
		var response httpresponse.RequestError
		json.Unmarshal(res.Body.Bytes(), &response)
		if response.Status != expectedStatus {
			t.Errorf(`expected status "%s" but received "%s"`, expectedStatus, response.Status)
		}
		if response.Error != expectedErrorDescription {
			t.Errorf(`expected description "%s" but received "%s"`, expectedErrorDescription, response.Error)
		}
	})

	t.Run("Error internal", func(t *testing.T) {
		t.Parallel()

		alias := "Q41Tqc"

		res := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/"+alias, nil)
		if err != nil {
			t.Fatalf("create new request: %v", err)
		}

		mockCacheURLGetter.On("GetURL", alias).Return("", nil)
		mockDBURLGetter.On("GetURL", alias).Return("", errors.New("internal"))

		mux.ServeHTTP(res, req)

		status := http.StatusInternalServerError
		if res.Code != status {
			t.Errorf("expected status code %d but received %d", status, res.Code)
		}

		expectedStatus := "Error"
		expectedErrorDescription := "getting url from storage: internal"
		var response httpresponse.RequestError
		json.Unmarshal(res.Body.Bytes(), &response)
		if response.Status != expectedStatus {
			t.Errorf(`expected status "%s" but received "%s"`, expectedStatus, response.Status)
		}
		if response.Error != expectedErrorDescription {
			t.Errorf(`expected description "%s" but received "%s"`, expectedErrorDescription, response.Error)
		}
	})

}
