package models

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"time"
)

type Schedule struct {
	Id      int64
	UserId  int64
	Weekday time.Weekday
	Start   time.Time
	End     time.Time
}

func (s Schedule) Validate(db *gorm.DB, errors *binding.Errors) error {
	if s.UserId == 0 {
		errors.Fields["UserId"] = "Invalid id"
	}

	if s.Weekday > time.Saturday {
		errors.Fields["Weekday"] = "Invalid weekday"
	}

	validTimes := true
	if s.Start.IsZero() {
		errors.Fields["Start"] = "Invalid time"
		validTimes = false
	}
	if s.End.IsZero() {
		errors.Fields["End"] = "Invalid time"
		validTimes = false
	}

	if validTimes && s.Start.After(s.End) {
		errors.Fields["Start"] = "Start time is after end time"
		validTimes = false
	}

	if validTimes {
		var count int
		query := db.Model(&Schedule{}).
			Where("weekday = ? AND (start >= ? OR end <= ?)", s.Weekday, s.End, s.Start).
			Count(&count)

		if query.Error != nil {
			return query.Error
		}

		if count > 0 {
			errors.Overall["Times"] = "Overlaps with existing schedule"
		}
	}

	return nil
}
