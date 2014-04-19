package api

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"ridesyncer/models"
)

const ApiTokenHeaderKey = "X-API-TOKEN"

func AuthenticateUser(db *gorm.DB) martini.Handler {
	return func(req *http.Request, c martini.Context) {
		token := req.Header.Get(ApiTokenHeaderKey)
		var authUser models.AuthUser
		if token != "" {
			var err error
			if authUser.User, err = models.GetUserByToken(db, token); err == nil {
				authUser.SetAuthenticated(true)
			}
		}
		c.Map(authUser)
	}
}

// Makes sure an authUser is authenticated and optionally verified
func NeedsAuth(checkVerified bool) martini.Handler {
	return func(res http.ResponseWriter, authUser models.AuthUser) {
		if !authUser.IsAuthenticated() || (checkVerified && !authUser.EmailVerified) {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
