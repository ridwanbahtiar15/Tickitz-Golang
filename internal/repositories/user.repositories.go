package repositories

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*sqlx.DB
}

type IUserRepository interface {
	RepositoryGetUserProfile(ID int) ([]models.UserProfileModel, error)
	RepositorySensitiveData(ID int) ([]models.UserProfileModel, error)
	RepositoryUpdateUser(userID int, body *models.UserProfileUpdateModel, url, hashedPassword string) (int64, error)
}

func InitializeUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) RepositoryGetUserProfile(ID int) ([]models.UserProfileModel, error) {
	data := []models.UserProfileModel{}
	query := `select u.id as "no",
	COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') AS "full_name",
	u.user_photo_profile as "user_photo",
	u.first_name as "firstname",
	u.last_name as "lastname",
	u.phone as "phone",
	u.email as "email",
	u.user_role as "user_role",
	u.created_at as "created_at"
	from users u
	where u.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositorySensitiveData(ID int) ([]models.UserProfileModel, error) {
	data := []models.UserProfileModel{}
	query := `select u.user_role as "user_role",
	u.otp as "otp",
	u.password_user as "password",
	u.activated as "activate"
	from users u
	where u.id = $1`
	err := r.Select(&data, query, ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func (r *UserRepository) RepositoryUpdateUser(userID int, body *models.UserProfileUpdateModel, url, hashedPassword string) (int64, error) {
	var conditional []string
	query := `
        UPDATE users
        SET `
	params := make(map[string]interface{})
	if url != "" {
		conditional = append(conditional, "user_photo_profile = :url")
		params["url"] = url
	}
	if body.First_Name != "" {
		conditional = append(conditional, "first_name = :firstname")
		params["firstname"] = body.First_Name
	}
	if body.Last_Name != "" {
		conditional = append(conditional, "last_name = :lastname")
		params["lastname"] = body.Last_Name
	}
	if body.Phone != nil {
		conditional = append(conditional, "phone = :phone")
		params["phone"] = body.Phone
	}
	if hashedPassword != "" {
		conditional = append(conditional, "password_user = :password")
		params["password"] = hashedPassword
	}
	if body.Role != "" {
		conditional = append(conditional, "user_role = :user_role")
		params["user_role"] = body.Role
	}
	if len(conditional) == 1 {
		query += conditional[0]
	}
	if len(conditional) > 1 {
		query += strings.Join(conditional, ", ")
	}
	params["Id"] = userID
	query += ` ,updated_at = NOW() WHERE id = :Id`
	// fmt.Println(query)
	result, err := r.NamedExec(query, params)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}
