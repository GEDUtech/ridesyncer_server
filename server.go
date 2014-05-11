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

func setupDb(m *martini.ClassicMartini) *gorm.DB {
	db, err := gorm.Open("mysql", os.Getenv("RIDESYNCER_DB_SOURCE"))

	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %s", err))
	}

	return &db
}

func setupEmail(m *martini.ClassicMartini) *email.Config {
	emailConfig, err := email.NewConfig(
		os.Getenv("RIDESYNCER_EMAIL_USERNAME"),
		os.Getenv("RIDESYNCER_EMAIL_PASSWORD"),
		os.Getenv("RIDESYNCER_EMAIL_HOST"),
		os.Getenv("RIDESYNCER_EMAIL_PORT"))

	if err != nil {
		panic(err)
	}

	return emailConfig
}

func main() {
	// Create martini
	m := martini.Classic()

	// Setup database connection
	db := setupDb(m)

	// Setup email configuration
	emailConfig := setupEmail(m)

	// Setup middleware
	m.Use(martini.Static("public"))
	m.Use(render.Renderer())
	m.Use(api.AuthenticateUser(db))

	// Create controllers
	pages := controllers.PagesController{}
	apiUsers := api.NewUsers(db, emailConfig)
	apiSchedules := api.NewSchedules(db)
	apiSyncs := api.NewSyncs(db)

	// Routing
	m.Get("/", pages.Home)

	m.Group("/api/users", func(r martini.Router) {
		r.Post("/login", apiUsers.Login)
		r.Post("/register", binding.Form(models.RegisterUser{}), apiUsers.Register)
		r.Post("/register_gcm", api.NeedsAuth(true), apiUsers.RegisterGcm)
		r.Post("/verify", api.NeedsAuth(false), apiUsers.Verify)
		r.Get("/search", api.NeedsAuth(true), apiUsers.Search)
	})

	m.Group("/api/schedules", func(r martini.Router) {
		r.Post("/add", api.NeedsAuth(true), apiSchedules.Add)
	})

	m.Group("/api/syncs", func(r martini.Router) {
		r.Post("/create", apiSyncs.Create)
	})

	// Start server
	m.Run()
}
