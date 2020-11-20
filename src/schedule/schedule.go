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
	Groups     []string
	OwnerID    string
	Season     string
	ScheduleID string
}

// AddSchedule creates an schedule for user in season and saves it to database
func AddSchedule(ownerID string, seasonID string) (Schedule, error) {
	schedule := Schedule{
		OwnerID:    ownerID,
		Groups:     []string{},
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

// AddGroup adds a new course to the schedule specified and saves the change to database
func (s *Schedule) AddGroup(groupID string) error {
	for _, group := range s.Groups {
		if group == groupID {
			return ErrAlreadyInGroup
		}
	}
	s.Groups = append(s.Groups, groupID)
	collection := database.DbClient.Database("test").Collection("schedules")
	filter := bson.M{
		"scheduleid": s.ScheduleID,
	}
	group, err := GetGroup(groupID)
	if err != nil {
		return err
	}
	owner, err := user.GetUser(s.OwnerID)
	if err != nil {
		return err
	}
	collection.FindOneAndReplace(context.TODO(), filter, *s)
	group.AddStudent(owner)
	return nil
}

// DeleteGroup a group from database
func DeleteGroup(groupID string) {
	collection := database.DbClient.Database("test").Collection("groups")
	filter := bson.M{
		"groupid": groupID,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
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

// DeleteSchedule removes schedule
func DeleteSchedule(scheduleID string) {
	collection := database.DbClient.Database("test").Collection("schedules")
	filter := bson.M{
		"scheduleid": scheduleID,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
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

func getSeasons() ([]Season, error) {
	collection := database.DbClient.Database("test").Collection("seasons")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	seasons := []Season{}
	for cur.Next(context.TODO()) {
		var elem Season
		cur.Decode(&elem)
		seasons = append(seasons, elem)
	}
	return seasons, nil
}

// GetSeason gets a season from the database
func GetSeason(ID string) (Season, error) {
	var season Season
	filter := bson.M{
		"id": ID,
	}
	err := database.DbClient.Database("test").Collection("seasons").FindOne(context.TODO(), filter).Decode(&season)
	if err != nil {
		return Season{}, err
	}
	return season, nil
}

// DeleteSeason removes season from database
func DeleteSeason(ID string) {
	collection := database.DbClient.Database("test").Collection("seasons")
	filter := bson.M{
		"id": ID,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
}
