package api

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"ridesyncer/models"
	"ridesyncer/utils"
)

const ApiTokenHeaderKey = "X-API-TOKEN"

func AuthenticateUser(db *gorm.DB) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		token := req.Header.Get(ApiTokenHeaderKey)
		var authUser models.AuthUser
		if token != "" {
			var err error
			authUser.User, err = models.GetUserByToken(db, token)

			if err == nil {
				authUser.SetAuthenticated(true)
			} else if err != gorm.RecordNotFound {
				utils.HttpError(res, http.StatusInternalServerError)
			}
		}
		c.Map(authUser)
	}
}

// Makes sure an authUser is authenticated and optionally verified
func NeedsAuth(checkVerified bool) martini.Handler {
	return func(res http.ResponseWriter, authUser models.AuthUser) {
		if !authUser.IsAuthenticated() || (checkVerified && !authUser.EmailVerified) {
			utils.HttpError(res, http.StatusUnauthorized)
		}
	}
}
