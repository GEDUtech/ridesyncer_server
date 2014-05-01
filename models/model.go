package models

import (
	"ridesyncer/utils"
)

func newValidation(errors *Errors) *utils.Validation {
	return &utils.Validation{func(field, message string) {
		errors.Fields[field] = message
	}}
}
