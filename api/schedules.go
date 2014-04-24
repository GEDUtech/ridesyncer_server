package api

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"ridesyncer/models"
)

type Schedules struct {
	db *gorm.DB
}

func NewSchedules(db *gorm.DB) Schedules {
	return Schedules{db}
}

func (this *Schedules) Add(req *http.Request, authUser models.AuthUser, render render.Render) {
	var schedule models.Schedule

	if decode(req, render, &schedule) != nil {
		return
	}

	schedule.Id = 0
	schedule.UserId = authUser.Id

	errors := binding.Errors{make(map[string]string), make(map[string]string)}
	if err := schedule.Validate(this.db, &errors); err != nil {
		render.Error(500)
		return
	}

	if errors.Count() > 0 {
		render.JSON(400, errors)
		return
	}

	if this.db.Save(&schedule).Error != nil {
		render.Error(500)
		return
	}
}
