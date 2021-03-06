package api

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"html/template"
	"net/http"
	"ridesyncer/auth"
	"ridesyncer/collections"
	"ridesyncer/models"
	"ridesyncer/net/email"
	"ridesyncer/utils"
)

var (
	emailTemplate *template.Template
)

type Users struct {
	db          *gorm.DB
	emailConfig *email.Config
}

func NewUsers(db *gorm.DB, emailConfig *email.Config) Users {
	return Users{db, emailConfig}
}

func (this *Users) Login(res http.ResponseWriter, req *http.Request, render render.Render) {
	var data struct{ Username, Password string }

	if decode(req, render, &data) != nil {
		return
	}

	user, err := models.GetUserByUsername(this.db, data.Username)
	if err != nil {
		switch err {
		case gorm.RecordNotFound:
			utils.HttpError(res, http.StatusUnauthorized)
		default:
			utils.HttpError(res, http.StatusInternalServerError)
		}
		return
	}

	if auth.NewBcryptHasher().Check(user.Password, data.Password) != nil {
		utils.HttpError(res, http.StatusUnauthorized)
		return
	}

	token, err := models.GenerateApiToken(this.db)
	if err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	user.Token = token
	if this.db.Model(user).UpdateColumn("token", token).Error != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	user.Password = ""

	if user.FetchSchedules(this.db) != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	if user.FetchSyncs(this.db) != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	render.JSON(http.StatusOK, user)
}

func (this *Users) Register(res http.ResponseWriter, req *http.Request, render render.Render) {
	var user models.RegisterUser
	if decode(req, render, &user) != nil {
		return
	}

	errors := models.NewErrors()
	if err := user.Validate(this.db, errors); err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	if errors.Count() > 0 {
		render.JSON(http.StatusBadRequest, map[string]*models.Errors{"errors": errors})
		return
	}

	if err := user.Geocode(); err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	if err := user.Register(this.db); err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	token, err := models.GenerateApiToken(this.db)
	if err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	user.Token = token
	if this.db.Model(user.User).UpdateColumn("token", token).Error != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}

	user.Password = ""
	render.JSON(http.StatusOK, user)

	go this.sendVerificationCode(user.User)
}

func (this *Users) RegisterGcm(res http.ResponseWriter, req *http.Request, authUser models.AuthUser, render render.Render) {
	var data struct{ GcmRegid string }
	if decode(req, render, &data) != nil {
		return
	}

	if this.db.Model(&authUser.User).Update("gcm_regid", data.GcmRegid).Error != nil {
		utils.HttpError(res, http.StatusInternalServerError)
	}
}

func (this *Users) Verify(res http.ResponseWriter, req *http.Request, authUser models.AuthUser, render render.Render) {
	if authUser.EmailVerified {
		utils.HttpError(res, 422)
		return
	}

	var data struct{ VerificationCode string }
	if decode(req, render, &data) != nil {
		return
	}

	if authUser.VerificationCode != data.VerificationCode {
		utils.HttpError(res, http.StatusBadRequest)
		return
	}

	query := this.db.Model(&authUser.User).Updates(map[string]interface{}{
		"email_verified":    true,
		"verification_code": "",
	})

	if query.Error != nil {
		utils.HttpError(res, http.StatusInternalServerError)
	}
}

func (this *Users) Search(res http.ResponseWriter, authUser models.AuthUser, render render.Render) {
	users, err := collections.FindMatches(this.db, authUser.User)
	if err != nil {
		utils.HttpError(res, http.StatusInternalServerError)
		return
	}
	render.JSON(http.StatusOK, map[string]interface{}{"results": users})
}

func (this *Users) sendVerificationCode(user models.User) error {
	var buffer bytes.Buffer
	viewData := map[string]string{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"V1":        user.VerificationCode[:3],
		"V2":        user.VerificationCode[3:6],
		"V3":        user.VerificationCode[6:9],
	}
	err := emailTemplate.Execute(&buffer, viewData)

	if err != nil {
		return err
	}

	return email.NewMailer(this.emailConfig).
		AddTo(user.Email).
		SetContentType("text/html").
		SetSubject("Verification Code").
		SetCharset("UTF-8").
		Send(buffer.String())
}

func init() {
	emailTemplate = template.Must(template.ParseFiles("templates/email/account_verification.tmpl"))
}
