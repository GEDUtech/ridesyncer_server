// Copyright 2014 GEDUtech. All rights reserved.

package controllers

import (
	"github.com/martini-contrib/render"
)

// Controller for Pages
type PagesController struct{}

// Renders the index (home) page
func (this *PagesController) Home(r render.Render) {
	r.HTML(200, "pages/home", nil)
}
