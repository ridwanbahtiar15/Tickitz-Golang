package repositories

import (
	"gilangrizaltin/Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (r *AuthRepositoryMock) RepositoryRegisterUser(body *models.AuthUserModel, hashedPassword string, otp int) error {
	args := r.Mock.Called(body, hashedPassword, otp)
	return args.Error(0)
}

func (r *AuthRepositoryMock) RepositorySelectPrivateData(email string) ([]models.UserProfileModel, error) {
	args := r.Mock.Called(email)
	return args.Get(0).([]models.UserProfileModel), args.Error(1)
}

func (r *AuthRepositoryMock) RepositoryActivateAccount(email string) error {
	args := r.Mock.Called(email)
	return args.Error(0)
}

func (r *AuthRepositoryMock) RepositoryLogOut(token string) error {
	args := r.Mock.Called(token)
	return args.Error(0)
}

func (r *AuthRepositoryMock) RepositoryIsTokenBlacklisted(email string) (bool, error) {
	args := r.Mock.Called(email)
	return args.Get(0).(bool), args.Error(1)
}
