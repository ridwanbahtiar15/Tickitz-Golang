package handlers

import (
	"encoding/json"
	"errors"
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var arm = repositories.AuthRepositoryMock{}
var handler = InitializeAuthHandler(&arm)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth", handler.Register)
	t.Run("Validation error in register", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin",
			Password: "1231",
		}
		arm.On("RepositoryRegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Input not valid", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal Server Register Error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositoryRegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success register", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositoryRegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		dataRegister := make(map[string]interface{})
		dataRegister["email"] = body.Email
		expectedMessage := helpers.NewResponse("Successfully register", dataRegister, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestActivate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth", handler.ActivateAccount)
	t.Run("Validation error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, nil)
		arm.On("RepositoryActivateAccount", mock.Anything).Return(nil)
		req := httptest.NewRequest("POST", "/auth?email=gilangmrizaltin&otp=821321", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Input not valid", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal error private data", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, errors.New("error"))
		arm.On("RepositoryActivateAccount", mock.Anything).Return(nil)
		req := httptest.NewRequest("POST", "/auth?email=gilangmrizaltin@gmail.com&otp=821321", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error in Private data", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, nil)
		arm.On("RepositoryActivateAccount", mock.Anything).Return(nil)
		req := httptest.NewRequest("POST", "/auth?email=gilangmrizaltin@gmail.com&otp=821321", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Account not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal server error activating account", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		// data := make([]models.UserProfileModel, 1)
		dataCorrect := []models.UserProfileModel{{Otp: 821321}}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return(dataCorrect, nil)
		arm.On("RepositoryActivateAccount", mock.Anything).Return(errors.New("error"))
		req := httptest.NewRequest("POST", "/auth?email=gilangmrizaltin@gmail.com&otp=821321", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error in activating account", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success activate", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		// data := make([]models.UserProfileModel, 1)
		dataCorrect := []models.UserProfileModel{{Otp: 821321}}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return(dataCorrect, nil)
		arm.On("RepositoryActivateAccount", mock.Anything).Return(nil)
		req := httptest.NewRequest("POST", "/auth?email=gilangmrizaltin@gmail.com&otp=821321", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Activate account success", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth", handler.Login)
	t.Run("Validation Error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Input not valid", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal Server Error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, errors.New("error"))
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error in Private data", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Account not found", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return([]models.UserProfileModel{}, nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Account not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Invalid password", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		dataCorrect := []models.UserProfileModel{{Activate: true, Password: "$argon2id$v=19$m=65536,t=3,p=2$Q+sjfQ8jZltu3KDHWneEqw$ZymDcHg3j083lGZelKd9HWyT6ICAy5Y3EcKloIUEBCs"}}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return(dataCorrect, nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Email or password is wrong", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Successfully login", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.AuthUserModel{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "12345",
		}
		dataCorrect := []models.UserProfileModel{{Activate: true, Password: "$argon2id$v=19$m=65536,t=3,p=2$Q+sjfQ8jZltu3KDHWneEqw$ZymDcHg3j083lGZelKd9HWyT6ICAy5Y3EcKloIUEBCs"}}
		arm.On("RepositorySelectPrivateData", mock.Anything).Return(dataCorrect, nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth", handler.Logout)
}
