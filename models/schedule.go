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
	Start   string
	End     string
}

func (s Schedule) Validate(db *gorm.DB, errors *binding.Errors) error {
	if s.UserId == 0 {
		errors.Fields["UserId"] = "Invalid id"
	}

	if s.Weekday > time.Saturday {
		errors.Fields["Weekday"] = "Invalid weekday"
	}

	validTimes := true

	start, _ := time.Parse("15:04:05", s.Start)
	end, _ := time.Parse("15:04:05", s.End)

	if start.IsZero() {
		errors.Fields["Start"] = "Invalid time"
		validTimes = false
	}
	if end.IsZero() {
		errors.Fields["End"] = "Invalid time"
		validTimes = false
	}

	if validTimes && start.After(end) {
		errors.Fields["Start"] = "Start time is after end time"
		validTimes = false
	}

	if validTimes && s.UserId != 0 {
		var count int
		query := db.Model(&Schedule{}).
			Where("user_id = ? AND weekday = ? AND (start >= ? OR end <= ?)", s.UserId, int(s.Weekday), s.Start, s.End).
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
