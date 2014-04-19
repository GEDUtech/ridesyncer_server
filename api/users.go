package api

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"html/template"
	"net/http"
	"net/smtp"
	"ridesyncer/auth"
	"ridesyncer/models"
	"ridesyncer/net/email"
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

func (this *Users) Login(req *http.Request, render render.Render) {
	user, err := models.GetUserByUsername(this.db, req.FormValue("username"))

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

	token, err := models.GenerateApiToken(this.db)
	if err != nil {
		render.Error(500) // Shouldn't happen
		return
	}

	user.Token = token
	if this.db.Model(user).UpdateColumn("token", token).Error != nil {
		render.Error(500) // Shouldn't happen
		return
	}

	user.Password = ""
	render.JSON(200, user)
}

func (this *Users) Register(user models.RegisterUser, validationErrors binding.Errors, render render.Render) {
	if err := user.Validate(this.db, &validationErrors); err != nil {
		render.JSON(500, err)
		return
	}

	if validationErrors.Count() > 0 {
		render.JSON(400, map[string]binding.Errors{"errors": validationErrors})
		return
	}

	if err := user.Register(this.db); err != nil {
		render.JSON(500, err)
		return
	}

	go this.sendVerificationCode(user.User)
}

func (this *Users) Verify(req *http.Request, authUser models.AuthUser, render render.Render) {
	if authUser.EmailVerified {
		render.Error(422)
		return
	}

	verificationCode := req.FormValue("VerificationCode")
	if authUser.VerificationCode != verificationCode {
		render.Error(400)
		return
	}

	query := this.db.Model(&authUser.User).Updates(map[string]interface{}{
		"email_verified":    true,
		"verification_code": "",
	})

	if query.Error != nil {
		render.Error(500)
	}
}

func (this *Users) sendVerificationCode(user models.User) error {
	subject := "Subject: Verification Code\n"
	mime := "MIME-version: 1.0;\n"
	contentType := `Content-Type: text/html; charset="UTF-8";`

	buffer := bytes.NewBufferString(subject + mime + contentType + "\n\n")
	viewData := map[string]string{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"V1":        user.VerificationCode[:3],
		"V2":        user.VerificationCode[3:6],
		"V3":        user.VerificationCode[6:9],
	}
	emailTemplate.Execute(buffer, viewData)

	return smtp.SendMail(this.emailConfig.Addr, this.emailConfig.Auth, "RideSyncer",
		[]string{user.Email}, buffer.Bytes())
}

func init() {
	emailTemplate = template.Must(template.ParseFiles("templates/email/account_verification.tmpl"))
}
