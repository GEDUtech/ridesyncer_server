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

// Makes sure a user attempting request exists
func TokenRequired(res http.ResponseWriter, req *http.Request, authUser models.AuthUser) {
	if !authUser.IsAuthenticated() {
		unauthorized(res)
	}
}

// Makes sure the user exists and is verified
func Authenticated(res http.ResponseWriter, authUser models.AuthUser) {
	if !authUser.IsAuthenticated() || !authUser.EmailVerified {
		unauthorized(res)
	}
}

func unauthorized(res http.ResponseWriter) {
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
