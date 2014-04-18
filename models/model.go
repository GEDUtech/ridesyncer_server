package models

import (
	"github.com/martini-contrib/binding"
	"ridesyncer/utils"
)

func newValidation(errors *binding.Errors) *utils.Validation {
	return &utils.Validation{func(field, message string) {
		errors.Fields[field] = message
	}}
}
