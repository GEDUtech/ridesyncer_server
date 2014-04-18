package api

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"net/http"
	"ridesyncer/auth"
	"ridesyncer/models"
)

type Users struct{}

func (this *Users) Login(req *http.Request, db gorm.DB, render render.Render) {
	user, err := models.GetUserByUsername(db, req.FormValue("username"))

	if err != nil {
		switch err {
		case gorm.RecordNotFound:
			render.Error(401) // Not authorized (invalid login)
		default:
			render.Error(500) // Shouldn't happen
		}
		return
	}

	if auth.NewBcryptHasher().Check(user.Password, req.FormValue("password")) != nil {
		render.Error(401) // Not authorized (invalid password)
		return
	}

	token, err := models.GenerateApiToken(db)
	if err != nil {
		render.Error(500) // Shouldn't happen
		return
	}

	user.Token = token
	if db.Model(user).UpdateColumn("token", token).Error != nil {
		render.Error(500) // Shouldn't happen
		return
	}

	user.Password = ""
	render.JSON(200, user)
}
