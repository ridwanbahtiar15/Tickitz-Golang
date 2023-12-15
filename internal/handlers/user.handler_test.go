package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var urm = repositories.UserRepositoryMock{}
var handlerUser = InitializeUserHandler(&urm)

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user", handlerUser.GetUserProfile)
	t.Run("Internal Server Error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		urm.On("RepositoryGetUserProfile", mock.Anything).Return([]models.UserProfileModel{}, errors.New("error"))
		req := httptest.NewRequest("GET", "/user", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		urm.On("RepositoryGetUserProfile", mock.Anything).Return([]models.UserProfileModel{}, nil)
		req := httptest.NewRequest("GET", "/user", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data user not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get profile", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		data := make([]models.UserProfileModel, 1)
		urm.On("RepositoryGetUserProfile", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("GET", "/user", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully Get Profile user", data, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestEditProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	var dataNotExecuted int64 = 0
	var dataExecuted int64 = 1
	r.PATCH("/user", handlerUser.UpdateProfileUser)
	t.Run("Validation Error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UserProfileUpdateModel{
			First_Name: "Gil@ng",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("firstname", body.First_Name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExecuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Input not valid", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Mismatch stored and new password", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UserProfileUpdateModel{
			New_Password:  "12343",
			Last_Password: "12412",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("last_password", body.Last_Password)
		_ = writer.WriteField("new_password", body.New_Password)
		writer.Close()
		dataCorrect := []models.UserProfileModel{{Password: "$argon2id$v=19$m=65536,t=3,p=2$Q+sjfQ8jZltu3KDHWneEqw$ZymDcHg3j083lGZelKd9HWyT6ICAy5Y3EcKloIUEBCs"}}
		urm.On("RepositorySensitiveData", mock.Anything).Return(dataCorrect, nil)
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataExecuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Last password doesnt match", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Internal Server Error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UserProfileUpdateModel{
			First_Name: "Gilang",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("firstname", body.First_Name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExecuted, errors.New("error"))
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("User not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UserProfileUpdateModel{
			First_Name: "Gilang",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("firstname", body.First_Name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExecuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("User not found", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Success update user", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UserProfileUpdateModel{
			First_Name: "Gilang",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("firstname", body.First_Name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataExecuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully update user", body, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
}
