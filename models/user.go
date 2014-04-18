package models

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id               int64  `form:"id"`
	Username         string `form:"username"`
	Password         string `form:"password"`
	FirstName        string `form:"first_name"`
	LastName         string `form:"last_name"`
	Email            string `form:"email"`
	EmailVerified    bool
	VerificationCode string
	Ride             string `form:"ride"`
	Token            string `form:"token"`
	authenticated    bool   `sql:"-" form:"-"`
	CreatedAt        time.Time
}

func GenerateApiToken(db gorm.DB) (string, error) {
	for {
		size := 20
		randomBytes := make([]byte, size)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return "", err
		}

		hexBytes := make([]byte, hex.EncodedLen(size))
		hex.Encode(hexBytes, randomBytes)
		token := string(hexBytes)

		var count int
		query := db.Model(&User{}).Where(&User{Token: token}).Count(&count)

		if query.Error != nil || count == 0 {
			return token, err
		}
	}
}

func (user *User) SetAuthenticated(authenticated bool) {
	user.authenticated = authenticated
}

func (user *User) IsAuthenticated() bool {
	return user.authenticated
}

func GetUserByToken(db gorm.DB, token string) (user User, err error) {
	err = db.Find(&user, User{Token: token}).Error
	return
}

func GetUserByUsername(db gorm.DB, username string) (user *User, err error) {
	user = new(User)
	err = db.Find(user, User{Username: username}).Error
	return
}
