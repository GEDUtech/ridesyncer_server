// Copyright 2014 GEDUtech. All rights reserved.

package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"ridesyncer/controllers"
)

func main() {
	// Create martini
	m := martini.Classic()

	// Setup middleware
	m.Use(martini.Static("public"))
	m.Use(render.Renderer())

	// Create controllers
	pages := controllers.PagesController{}

	// Routing
	m.Get("/", pages.Home)

	// Start server
	m.Run()
}
