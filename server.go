// Copyright 2014 GEDUtech. All rights reserved.

package main

import (
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"os"
	"ridesyncer/api"
	"ridesyncer/controllers"
	"ridesyncer/models"
	"ridesyncer/net/email"
)

func setupDb(m *martini.ClassicMartini) {
	db, err := gorm.Open("mysql", os.Getenv("RIDESYNCER_DB_SOURCE"))

	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %s", err))
	}

	m.Map(db)
}

func setupEmail(m *martini.ClassicMartini) {
	emailConfig, err := email.NewConfig(
		os.Getenv("RIDESYNCER_EMAIL_USERNAME"),
		os.Getenv("RIDESYNCER_EMAIL_PASSWORD"),
		os.Getenv("RIDESYNCER_EMAIL_HOST"),
		os.Getenv("RIDESYNCER_EMAIL_PORT"))

	if err != nil {
		panic(err)
	}

	m.Map(emailConfig)
}

func main() {
	// Create martini
	m := martini.Classic()

	// Setup database connection
	setupDb(m)

	// Setup email configuration
	setupEmail(m)

	// Setup middleware
	m.Use(martini.Static("public"))
	m.Use(render.Renderer())
	m.Use(api.AuthenticateUser)

	// Create controllers
	pages := controllers.PagesController{}
	apiUsers := api.Users{}

	// Routing
	m.Get("/", pages.Home)

	m.Group("/api/users", func(r martini.Router) {
		r.Post("/login", apiUsers.Login)
		r.Post("/register", binding.Form(models.RegisterUser{}), apiUsers.Register)
	})

	// Start server
	m.Run()
}
