package models

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
)

type RegisterUser struct {
	User

	RepeatPassword string `sql:"-" form:"repeat_password"`
}

func (registerUser *RegisterUser) Validate(db gorm.DB, errors *binding.Errors) error {
	validation := newValidation(errors)

	if validation.NotEmpty("repeat_password", registerUser.RepeatPassword) {
		validation.Match("repeat_password", registerUser.Password, registerUser.RepeatPassword, "Passwords do not match")
	}

	return registerUser.User.Validate(db, errors)
}
