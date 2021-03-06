package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	geo "github.com/kellydunn/golang-geo"
	mrand "math/rand"
	"ridesyncer/auth"
	"strconv"
	"time"
)

type User struct {
	Id               int64
	Username         string
	Password         string `json:",omitempty"`
	FirstName        string
	LastName         string
	Email            string
	EmailVerified    bool
	VerificationCode string `json:",omitempty"`
	Ride             string
	Address          string  `json:",omitempty"`
	City             string  `json:",omitempty"`
	State            string  `json:",omitempty"`
	Zip              string  `json:",omitempty"`
	Token            string  `json:",omitempty"`
	Buffer           int     `json:",omitempty"`
	Lat              float64 `json:"-"`
	Lng              float64 `json:"-"`
	GcmRegid         sql.NullString
	CreatedAt        time.Time
	Distance         float64 `sql:"-"`

	Schedules []Schedule `sql:"-"`
	Syncs     []Sync     `sql:"-"`
}

func (user *User) Geocode() error {
	address := fmt.Sprintf("%s %s, %s %s", user.Address, user.City, user.State, user.Zip)

	geocoder := geo.GoogleGeocoder{}
	point, err := geocoder.Geocode(address)

	if err != nil {
		return err
	}

	user.Lat = point.Lat()
	user.Lng = point.Lng()

	return nil
}

func (user *User) FetchSchedules(db *gorm.DB) error {
	return db.Model(user).Related(&user.Schedules).Error
}

func (user *User) FetchSyncs(db *gorm.DB) error {
	query := db.Table("syncs AS Sync").Select("Sync.id, Sync.weekday, Sync.created_at").
		Joins("LEFT JOIN sync_users AS SyncUser ON (SyncUser.sync_id = Sync.id)").
		Where("SyncUser.user_id = ?", user.Id).
		Find(&user.Syncs)

	if query.Error != nil {
		return query.Error
	}

	if len(user.Syncs) == 0 {
		return nil
	}

	ids := make([]int64, len(user.Syncs))
	syncsMap := map[int64]*Sync{}
	for idx, sync := range user.Syncs {
		ids[idx] = sync.Id
		syncsMap[sync.Id] = &user.Syncs[idx]
	}

	syncUsers := []SyncUser{}
	query = db.Model(&SyncUser{}).Where("sync_id IN (?)", ids).Find(&syncUsers)
	if query.Error != nil {
		return query.Error
	}

	for _, syncUser := range syncUsers {
		if err := db.Model(&syncUser).Select("id, username, email, email_verified, first_name, last_name, ride").Related(&syncUser.User).Error; err != nil {
			return err
		}
		if err := db.Model(&syncUser.User).Related(&syncUser.User.Schedules).Error; err != nil {
			return err
		}
		syncsMap[syncUser.SyncId].SyncUsers = append(syncsMap[syncUser.SyncId].SyncUsers, syncUser)
	}

	return nil
}

func (user *User) ValidateUniqueUsername(db *gorm.DB, errors *Errors) error {
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

func (user *User) ValidateUniqueEmail(db *gorm.DB, errors *Errors) error {
	var count int
	query := db.Model(User{}).Where(&User{Email: user.Email}).Count(&count)
	if query.Error != nil {
		return query.Error
	}

	if count > 0 {
		errors.Fields["email"] = "Email already taken"
	}

	return nil
}

func (user *User) Validate(db *gorm.DB, errors *Errors) error {
	validation := newValidation(errors)

	if validation.Between("username", user.Username, 6, 16) {
		if err := user.ValidateUniqueUsername(db, errors); err != nil {
			return err
		}
	}

	validation.Between("password", user.Password, 6, 16)
	validation.NotEmpty("first_name", user.FirstName)
	validation.NotEmpty("last_name", user.LastName)

	if validation.NotEmpty("email", user.Email) && validation.Email("email", user.Email) {
		if err := user.ValidateUniqueEmail(db, errors); err != nil {
			return err
		}

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
