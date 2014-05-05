package models

import (
	"time"
)

type SyncUser struct {
	Id        int64
	UserId    int64
	SyncId    int64
	Status    int
	Order     int
	CreatedAt time.Time

	User User `sql:"-"`
}
