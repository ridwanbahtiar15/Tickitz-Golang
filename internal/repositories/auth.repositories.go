package repositories

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	*sqlx.DB
}

type IAuthRepository interface {
	RepositoryRegisterUser(body *models.AuthUserModel, hashedPassword string, otp int) error
	RepositorySelectPrivateData(email string) ([]models.UserProfileModel, error)
	RepositoryActivateAccount(email string) error
	RepositoryLogOut(token string) error
	RepositoryIsTokenBlacklisted(token string) (bool, error)
}

func InitializeAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) RepositoryRegisterUser(body *models.AuthUserModel, hashedPwd string, otp int) error {
	params := make(map[string]interface{})
	params["email"] = body.Email
	params["password"] = hashedPwd
	params["otp"] = otp
	query := `insert into users (email, password_user, otp, user_role) values (:email, :password, :otp, 'Normal User')`
	_, err := r.NamedExec(query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositorySelectPrivateData(email string) ([]models.UserProfileModel, error) {
	data := []models.UserProfileModel{}
	query := `select u.id as "no",
	u.user_photo_profile as "user_photo",
	u.first_name as "firstname",
	u.last_name as "lastname",
	u.password_user as "password",
	u.user_role as "user_role",
	u.otp as "otp",
	u.activated as "activate"
	from users u
	where u.email = $1`
	// values := []any{body.Email}
	err := r.Select(&data, query, email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func (r *AuthRepository) RepositoryActivateAccount(email string) error {
	params := make(map[string]interface{})
	params["email"] = email
	query := `update users set activated = true where email = :email`
	_, err := r.NamedExec(query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositoryLogOut(token string) error {
	query := `insert into jwt_blacklist (jwt_code) values ($1)`
	values := []any{token}
	_, err := r.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) RepositoryIsTokenBlacklisted(token string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM jwt_blacklist WHERE jwt_code = $1`
	err := r.Get(&count, query, token)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
