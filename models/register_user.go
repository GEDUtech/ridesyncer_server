package models

import (
	"github.com/jinzhu/gorm"
)

type RegisterUser struct {
	User

	RepeatPassword string `form:"RepeatPassword"`
}

func (registerUser *RegisterUser) Validate(db *gorm.DB, errors *Errors) error {
	validation := newValidation(errors)

	if validation.NotEmpty("repeat_password", registerUser.RepeatPassword) {
		validation.Match("repeat_password", registerUser.Password, registerUser.RepeatPassword, "Passwords do not match")
	}

	return registerUser.User.Validate(db, errors)
}
