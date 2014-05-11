package api

import (
	"database/sql"
	"fmt"
	"github.com/alexjlockwood/gcm"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"net/http"
	"ridesyncer/models"
	"ridesyncer/utils"
)

type Syncs struct {
	db *gorm.DB
}

func NewSyncs(db *gorm.DB) *Syncs {
	return &Syncs{db}
}

func (this *Syncs) Create(req *http.Request, authUser models.AuthUser, render render.Render) {
	syncs := []models.Sync{}
	if decode(req, render, &syncs) != nil {
		return
	}

	tx := this.db.Begin()
	userIds := utils.Set{}
	for _, sync := range syncs {
		if q := tx.Save(&sync); q.Error != nil {
			tx.Rollback()
			render.Error(500)
			return
		}

		for _, syncUser := range sync.SyncUsers {
			syncUser.SyncId = sync.Id
			if syncUser.UserId != authUser.Id {
				userIds.Insert(syncUser.UserId)
			}
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

	go this.sendPushNotifications(userIds.ToSlice())

	if authUser.FetchSyncs(this.db) != nil {
		render.Error(500)
		return
	}

	render.JSON(http.StatusOK, map[string]interface{}{"results": authUser.Syncs})
}

func (this *Syncs) sendPushNotifications(userIds []interface{}) {
	gcmIds := []sql.NullString{}
	if this.db.Model(&models.User{}).Debug().Where("id IN (?) AND gcm_regid IS NOT NULL", userIds).Pluck("gcm_regid", &gcmIds).Error == nil {
		regIds := []string{}
		for _, id := range gcmIds {
			regIds = append(regIds, id.String)
		}
		msg := &gcm.Message{RegistrationIDs: regIds, CollapseKey: "sync_request"}

		sender := &gcm.Sender{ApiKey: GOOGLE_API_KEY}
		response, err := sender.Send(msg, 2)
		if err != nil {
			fmt.Println("Failed to send message:", err)
			return
		}
		for _, itm := range response.Results {
			fmt.Println(itm.Error)
		}
	}
}
