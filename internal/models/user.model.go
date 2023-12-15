package models

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
	Photo_Profile interface{} `db:"user_photo" json:"user_photo" valid:"-"`
	First_Name    *string     `db:"firstname" form:"firstname" json:"firstname" valid:"-"`
	Last_Name     *string     `db:"lastname" form:"lastname" json:"lastname" valid:"-"`
	Email         string      `db:"email" form:"email" json:"email" valid:"-"`
	Password      string      `db:"password" form:"password" json:"password" valid:"-"`
	Role          string      `db:"user_role" form:"user_role" json:"user_role" valid:"-"`
	Otp           int         `db:"otp" form:"otp" json:"otp" valid:"-"`
	Points        int         `db:"points" form:"points" json:"points" valid:"-"`
	Activate      bool        `db:"activate" form:"activate" json:"activate" valid:"-"`
}
