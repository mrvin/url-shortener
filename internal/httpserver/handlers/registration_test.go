package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrvin/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/url-shortener/pkg/http/response"
	"github.com/stretchr/testify/mock"
)

type MockUserCreator struct {
	mock.Mock
}

func (m *MockUserCreator) CreateUser(_ context.Context, user *storage.User) error {
	args := m.Called(user.Name, user.Role)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		TestName                 string
		Username                 string
		Password                 string
		Role                     string
		StatusCode               int
		Error                    error
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:                 "Success smoke test",
			Username:                 "Bob",
			Password:                 "qwerty",
			Role:                     "user",
			StatusCode:               http.StatusCreated,
			Error:                    nil,
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:                 "Error username is short",
			Username:                 "b",
			Password:                 "qwerty",
			Role:                     "user",
			StatusCode:               http.StatusBadRequest,
			Error:                    nil,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "invalid request: tag: min value: b",
		},
		{
			TestName:                 "Error password is short",
			Username:                 "Bob",
			Password:                 "qwe",
			Role:                     "user",
			StatusCode:               http.StatusBadRequest,
			Error:                    nil,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "invalid request: tag: min value: qwe",
		},
		{
			TestName:                 "Error user already exists",
			Username:                 "Alice",
			Password:                 "qwerty",
			Role:                     "user",
			StatusCode:               http.StatusConflict,
			Error:                    storage.ErrUserExists,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "saving user to storage: user already exists",
		},
		{
			TestName:                 "Error internal",
			Username:                 "Jimmy",
			Password:                 "qwerty",
			Role:                     "user",
			StatusCode:               http.StatusInternalServerError,
			Error:                    errors.New("internal"),
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "saving user to storage: internal",
		},
	}

	mockCreator := new(MockUserCreator)
	handler := ErrorHandler("Registration user", NewRegistration(mockCreator))
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			dataRequest, err := json.Marshal(RequestRegistration{Username: test.Username, Password: test.Password})
			if err != nil {
				t.Fatalf("cant marshal json: %v", err)
			}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/users", bytes.NewReader(dataRequest))
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}

			mockCreator.On("CreateUser", test.Username, test.Role).Return(test.Error)

			handler.ServeHTTP(res, req)

			if res.Code != test.StatusCode {
				t.Errorf("expected status code %d but received %d", test.StatusCode, res.Code)
			}
			if test.StatusCode == http.StatusCreated {
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
