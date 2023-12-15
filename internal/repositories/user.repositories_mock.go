package repositories

import (
	"gilangrizaltin/Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) RepositoryGetUserProfile(ID int) ([]models.UserProfileModel, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).([]models.UserProfileModel), args.Error(1)
}

func (r *UserRepositoryMock) RepositorySensitiveData(ID int) ([]models.UserProfileModel, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).([]models.UserProfileModel), args.Error(1)
}

func (r *UserRepositoryMock) RepositoryUpdateUser(userID int, body *models.UserProfileUpdateModel, url, hashedPassword string) (int64, error) {
	args := r.Mock.Called(userID, body, url, hashedPassword)
	return args.Get(0).(int64), args.Error(1)
}
