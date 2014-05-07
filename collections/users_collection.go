package collections

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"ridesyncer/models"
	"strings"
	"time"
)

func FindMatches(db *gorm.DB, user models.User) ([]models.User, error) {
	if err := user.FetchSchedules(db); err != nil {
		return nil, err
	}

	if len(user.Schedules) == 0 {
		return []models.User{}, nil
	}

	params := []interface{}{}
	conditions := []string{}
	duration := time.Duration(user.Buffer) * time.Minute
	for _, schedule := range user.Schedules {
		start, _ := time.Parse("15:04:05", schedule.Start)
		end, _ := time.Parse("15:04:05", schedule.End)

		cond := `(
			Schedule.weekday = ? AND
			(Schedule.start = ? OR Schedule.start = ? OR SUBDATE(Schedule.start, INTERVAL User.buffer MINUTE) = ?) AND
			(Schedule.end = ? OR Schedule.end = ? OR ADDDATE(Schedule.end, INTERVAL User.buffer MINUTE) = ?)
		)`
		conditions = append(conditions, cond)
		params = append(params, int(schedule.Weekday),
			schedule.Start, start.Add(-duration).Format("15:04:05"), schedule.Start,
			schedule.End, end.Add(duration).Format("15:04:05"), schedule.End)
	}
	conditionsString := fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
	query := db.Where(conditionsString, params...)

	distance := fmt.Sprintf(`
		(3959 * acos(cos(radians(%f))
			* cos(radians(User.lat))
			* cos(radians(User.lng)
			- radians(%f))
			+ sin(radians(%f)) * sin(radians(User.lat)))
		) AS distance`,
		user.Lat, user.Lng, user.Lat)

	users := []models.User{}
	usersQuery := query.Model(models.User{}).
		Table("users AS User").
		Select("User.id, User.username, User.city, User.state, "+distance).
		Joins("LEFT JOIN schedules AS Schedule ON (User.id = Schedule.user_id)").
		Where("Schedule.user_id != ?", user.Id).
		Where("SUBSTRING_INDEX(User.email, '@', -1) = ?", strings.Split(user.Email, "@")[1]).
		Having("distance < 100").
		Group("User.id").
		Order("distance").
		Limit(20).
		Find(&users)

	if usersQuery.Error != nil {
		return nil, usersQuery.Error
	}

	if len(users) == 0 {
		return []models.User{}, nil
	}

	userIds := []int64{}
	usersMap := map[int64]*models.User{}
	for idx, user := range users {
		user.Schedules = []models.Schedule{}
		userIds = append(userIds, user.Id)
		usersMap[user.Id] = &users[idx]
	}

	var schedules []models.Schedule
	schedulesQuery := query.Model(models.Schedule{}).
		Table("schedules AS Schedule").
		Select("Schedule.*").
		Joins("LEFT JOIN users AS User ON (Schedule.user_id = User.id)").
		Where("Schedule.user_id IN (?)", userIds).
		Order("Schedule.user_id").
		Find(&schedules)

	if schedulesQuery.Error != nil {
		return nil, schedulesQuery.Error
	}

	for _, schedule := range schedules {
		usersMap[schedule.UserId].Schedules = append(usersMap[schedule.UserId].Schedules, schedule)
	}

	return users, nil
}
