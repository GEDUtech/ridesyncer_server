package models

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	mrand "math/rand"
	"ridesyncer/auth"
	"strconv"
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
	CreatedAt        time.Time
}

func (user *User) ValidateUniqueUsername(db *gorm.DB, errors *binding.Errors) error {
	var count int
	query := db.Model(User{}).Where(&User{Username: user.Username}).Count(&count)
	if query.Error != nil {
		return query.Error
	}

	if count > 0 {
		errors.Fields["username"] = "Username already taken"
	}

	return nil
}

func (user *User) Validate(db *gorm.DB, errors *binding.Errors) error {
	validation := newValidation(errors)

	if validation.Between("username", user.Username, 6, 16) {
		if err := user.ValidateUniqueUsername(db, errors); err != nil {
			return err
		}
	}

	validation.Between("password", user.Password, 6, 16)
	validation.NotEmpty("first_name", user.FirstName)
	validation.NotEmpty("last_name", user.LastName)

	if validation.NotEmpty("email", user.Email) {
		validation.Email("email", user.Email)
	}
	validation.NotEmpty("ride", user.Ride)

	return nil
}

func (user *User) Register(db *gorm.DB) error {
	var err error
	user.VerificationCode, err = GenerateVerificationCode(db)
	if err != nil {
		return err
	}

	user.Password = auth.NewBcryptHasher().Hash(user.Password)
	user.CreatedAt = time.Now()

	return db.Save(user).Error
}

func GetUserByToken(db *gorm.DB, token string) (user User, err error) {
	err = db.Find(&user, User{Token: token}).Error
	return
}

func GetUserByUsername(db *gorm.DB, username string) (user *User, err error) {
	user = new(User)
	err = db.Find(user, User{Username: username}).Error
	return
}

func GenerateApiToken(db *gorm.DB) (string, error) {
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

func GenerateVerificationCode(db *gorm.DB) (string, error) {
	mrand.Seed(time.Now().UnixNano())
	for {
		min := 100000000
		max := 999999999
		code := strconv.Itoa(mrand.Int()%(max-min) + min)

		var count int
		query := db.Model(&User{}).Where(&User{VerificationCode: code}).Count(&count)

		if query.Error != nil || count == 0 {
			return code, query.Error
		}
	}
}
