package schedule

import (
	"context"
	"wilhelmiina/database"
	"wilhelmiina/user"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Season represents one season of school where the schedule stays same
type Season struct {
	Name  string
	Start int64
	End   int64
	ID    string
}

// Schedule represents some users schedule
type Schedule struct {
	Groups     []Group
	OwnerID    string
	Season     string
	ScheduleID string
}

// AddSchedule creates an schedule for user in season and saves it to database
func AddSchedule(ownerID string, seasonID string) (Schedule, error) {
	schedule := Schedule{
		OwnerID:    ownerID,
		Groups:     []Group{},
		Season:     seasonID,
		ScheduleID: uuid.New().String(),
	}
	collection := database.DbClient.Database("test").Collection("schedules")
	_, err := collection.InsertOne(context.TODO(), schedule)
	if err != nil {
		return Schedule{}, err
	}

	// Add the schedule to the user that has it
	user, err := user.GetUser(ownerID)
	if err != nil {
		return Schedule{}, err
	}
	user.AddSchedule(schedule.ScheduleID)
	return schedule, nil
}

// AddCourseToSchedule adds a new course to the schedule specified
func AddCourseToSchedule(owner string, season int, course Course) error {
	// TODO Add course to schedule
	return nil
}

// GetScheduleForUser gets the schedule of an user in specific season
func GetScheduleForUser(ownerID string, seasonID string) (Schedule, error) {
	user, err := user.GetUser(ownerID)
	if err != nil {
		return Schedule{}, err
	}
	var schedule Schedule
	found := false
	for _, id := range user.ScheduleIDs {
		s, err := GetSchedule(id)
		if err != nil {
			return Schedule{}, err // Return error if something goes wrong in getting schedule
		}
		if s.Season == seasonID {
			found = true
			schedule = s
			break
		}
	}
	if !found {
		return Schedule{}, ErrScheduleNotFound
	}
	return schedule, nil
}

// GetSchedule gets schedule by id
func GetSchedule(scheduleID string) (Schedule, error) {
	var schedule Schedule
	filter := bson.M{
		"scheduleid": scheduleID,
	}
	err := database.DbClient.Database("test").Collection("schedules").FindOne(context.TODO(), filter).Decode(&schedule)
	if err != nil {
		return Schedule{}, err
	}
	return schedule, nil
}

// AddSeason adds season to database and returns it
func AddSeason(name string, start int64, end int64) (Season, error) {
	seasonID := uuid.New().String()

	season := Season{
		Name:  name,
		Start: start,
		End:   end,
		ID:    seasonID,
	}
	collection := database.DbClient.Database("test").Collection("seasons")
	_, err := collection.InsertOne(context.TODO(), season)
	if err != nil {
		return Season{}, err
	}
	return season, nil
}

func getSeasons() {
	// TODO get all seasons
}

func getSeasonWithIndex(index int) {
	// TODO get specific season
}
