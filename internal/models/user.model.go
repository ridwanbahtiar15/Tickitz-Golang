package models

import "time"

type AuthUserModel struct {
	Email    string `db:"email" form:"email" json:"email" valid:"email, required"`
	Password string `db:"password" form:"password" json:"" valid:"required"`
}

type AuthUserActivateModel struct {
	Email string `db:"email" form:"email" json:"email" valid:"email, required"`
	Otp   int    `db:"otp" form:"otp" json:"otp" valid:"numeric,required"`
}

type UserProfileModel struct {
	Id            int         `db:"no" valid:"-"`
	Photo_Profile interface{} `db:"user_photo" json:"user_photo"`
	Full_name     *string     `db:"full_name" json:"full_name"`
	First_Name    *string     `db:"firstname" form:"firstname" json:"firstname" valid:"matches(^[a-zA-Z ]+$), optional"`
	Last_Name     *string     `db:"lastname" form:"lastname" json:"lastname" valid:"matches(^[a-zA-Z ]+$), optional"`
	Email         string      `db:"email" form:"email" json:"email" valid:"optional"`
	Password      string      `db:"password" form:"password" json:"password" valid:"optional"`
	Role          string      `db:"user_role" form:"user_role" json:"user_role" valid:"in(Admin|Normal User),optional"`
	Otp           int         `db:"otp" form:"otp" json:"otp"`
	Points        int         `db:"points" form:"points" json:"points" valid:"optional"`
	Activate      bool        `db:"activate" form:"activate" json:"activate"`
	Phone         *string     `db:"phone" form:"phone" json:"phone" valid:"numeric, optional"`
	Created_at    *time.Time  `db:"created_at" json:"created_at"`
}

type UserProfileUpdateModel struct {
	Photo_Profile interface{} `db:"user_photo" json:"user_photo"`
	First_Name    string      `db:"firstname" form:"firstname" json:"firstname" valid:"matches(^[a-zA-Z ]+$), optional"`
	Last_Name     string      `db:"lastname" form:"lastname" json:"lastname" valid:"matches(^[a-zA-Z ]+$), optional"`
	Email         string      `db:"email" form:"email" json:"email" valid:"optional"`
	Last_Password string      `db:"last_password" form:"last_password" json:"last_password" valid:"optional"`
	New_Password  string      `db:"new_password" form:"new_password" json:"new_password" valid:"optional"`
	Role          string      `db:"user_role" form:"user_role" json:"user_role" valid:"in(Admin|Normal User),optional"`
	Points        int         `db:"points" form:"points" json:"points" valid:"optional"`
	Phone         *string     `db:"phone" form:"phone" json:"phone" valid:"numeric, optional"`
}
