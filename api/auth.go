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
		var user models.User
		if token != "" {
			var err error
			if user, err = models.GetUserByToken(db, token); err == nil {
				if user.EmailVerified {
					user.SetAuthenticated(true)
				}
			}
		}
		c.Map(user)
	}
}

// Makes sure a user attempting request exists
func TokenRequired(res http.ResponseWriter, req *http.Request, user models.User) {
	if user.Id == 0 {
		unauthorized(res)
	}
}

// Makes sure the user exists and is verified
func Authenticated(res http.ResponseWriter, user models.User) {
	if !user.IsAuthenticated() {
		unauthorized(res)
	}
}

func unauthorized(res http.ResponseWriter) {
	http.Error(res, "Not Authorized", http.StatusUnauthorized)
}
