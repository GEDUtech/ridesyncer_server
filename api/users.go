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

func (this *Users) Register(emailConfig *email.Config, db gorm.DB, user models.RegisterUser, validationErrors binding.Errors, render render.Render) {
	if err := user.Validate(db, &validationErrors); err != nil {
		render.JSON(500, err)
		return
	}

	if validationErrors.Count() > 0 {
		render.JSON(400, map[string]binding.Errors{"errors": validationErrors})
		return
	}

	if err := user.Register(db); err != nil {
		render.JSON(500, err)
		return
	}

	go sendVerificationCode(emailConfig, user.User)
}

func sendVerificationCode(emailConfig *email.Config, user models.User) error {
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

	return smtp.SendMail(emailConfig.Addr, emailConfig.Auth, "RideSyncer", []string{user.Email}, buffer.Bytes())
}

func init() {
	emailTemplate = template.Must(template.ParseFiles("templates/email/account_verification.tmpl"))
}
