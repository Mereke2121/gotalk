package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/gotalk/models"
	services "github.com/gotalk/pkg/services"
	mock_services "github.com/gotalk/pkg/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthorization, user *models.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "testname","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: models.User{
				UserName: "testname",
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_services.MockAuthorization, user *models.User) {
				s.EXPECT().AddUser(user).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `Authorized successfully`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)
			testCase.mockBehavior(auth, &testCase.inputUser)

			logger, _ := zap.NewProduction()
			service := &services.Service{Authorization: auth}
			handler := NewHandler(service, logger)

			// test server
			mux := chi.NewRouter()
			mux.Post("/sign-up", handler.SignUp)

			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// perform request
			mux.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
