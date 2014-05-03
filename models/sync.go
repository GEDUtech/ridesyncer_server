package models

import (
	"time"
)

type Sync struct {
	Id        int64
	Weekday   int
	CreatedAt time.Time

	SyncUsers []SyncUser `sql:"-"`
}
