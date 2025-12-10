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

type MockUserGetter struct {
	mock.Mock
}

func (m *MockUserGetter) GetUser(_ context.Context, name string) (*storage.User, error) {
	args := m.Called(name)
	return args.Get(0).(*storage.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	tests := []struct {
		TestName                 string
		Username                 string
		Password                 string
		StatusCode               int
		Error                    error
		User                     *storage.User
		ExpectedStatus           string
		ExpectedErrorDescription string
	}{
		{
			TestName:   "Success smoke test",
			Username:   "Bob",
			Password:   "qwerty",
			StatusCode: http.StatusOK,
			Error:      nil,
			User: &storage.User{
				Name:         "Bob",
				HashPassword: "$2a$10$WW6Hn.HPGq65LeLsk..b5O9.k4kQpgVWq7LeDwC9GGW9txOkrdybG",
				Role:         "user",
			},
			ExpectedStatus:           "OK",
			ExpectedErrorDescription: "",
		},
		{
			TestName:   "User is not authorized",
			Username:   "Alice",
			Password:   "password",
			StatusCode: http.StatusUnauthorized,
			Error:      nil,
			User: &storage.User{
				Name:         "Alice",
				HashPassword: "$2a$10$WW6Hn.HPGq65LeLsk..b5O9.k4kQpgVWq7LeDwC9GGW9txOkrdybG",
				Role:         "user",
			},
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "compare hash and password: crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			TestName:                 "Error internal",
			Username:                 "Jimmy",
			Password:                 "qwerty",
			StatusCode:               http.StatusInternalServerError,
			Error:                    errors.New("internal"),
			User:                     nil,
			ExpectedStatus:           "Error",
			ExpectedErrorDescription: "getting user from storage: internal",
		},
	}

	mockGetter := new(MockUserGetter)
	handler := ErrorHandler("Login user", NewLogin(mockGetter))
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			dataRequest, err := json.Marshal(RequestRegistration{Username: test.Username, Password: test.Password})
			if err != nil {
				t.Fatalf("cant marshal json: %v", err)
			}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/users/login", bytes.NewReader(dataRequest))
			if err != nil {
				t.Fatalf("cant create new request: %v", err)
			}

			mockGetter.On("GetUser", test.Username).Return(test.User, test.Error)

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
