package api

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"net/http"
	"ridesyncer/models"
)

type Syncs struct {
	db *gorm.DB
}

func NewSyncs(db *gorm.DB) *Syncs {
	return &Syncs{db}
}

func (this *Syncs) Create(req *http.Request, render render.Render) {
	syncs := []models.Sync{}
	if decode(req, render, &syncs) != nil {
		return
	}

	tx := this.db.Begin()
	for _, sync := range syncs {
		if q := tx.Save(&sync); q.Error != nil {
			tx.Rollback()
			render.Error(500)
			return
		}

		for _, syncUser := range sync.SyncUsers {
			syncUser.SyncId = sync.Id
			if q := tx.Save(&syncUser); q.Error != nil {
				tx.Rollback()
				render.Error(500)
				return
			}
		}
	}

	if q := tx.Commit(); q.Error != nil {
		tx.Rollback()
		render.Error(500)
		return
	}

	render.JSON(http.StatusOK, syncs)
}
