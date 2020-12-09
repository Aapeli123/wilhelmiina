package schedule

import (
	"context"

	"github.com/Aapeli123/wilhelmiina/database"
	"github.com/Aapeli123/wilhelmiina/user"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	// Check if schedule owner exists
	_, err := user.GetUser(ownerID)
	if err != nil {
		return Schedule{}, err
	}
	// Insert schedule to database
	collection := database.DbClient.Database("test").Collection("schedules")
	_, err = collection.InsertOne(context.TODO(), schedule)
	if err != nil {
		return Schedule{}, err
	}
	return schedule, nil
}

// RemoveSchedule deletes the specified schedule from database
func RemoveSchedule(scheduleID string) {
	collection := database.DbClient.Database("test").Collection("schedules")
	collection.DeleteOne(context.TODO(), bson.M{
		"scheduleid": scheduleID,
	})
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
	collection := database.DbClient.Database("test").Collection("schedules")
	var schedule Schedule
	err := collection.FindOne(context.TODO(), bson.M{
		"ownerid": ownerID, "season": seasonID,
	}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Schedule{}, ErrScheduleNotFound
		}
		return Schedule{}, err
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

func GetSeasons() ([]Season, error) {
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

// DeleteSeason removes season from database, it also should remove all schedules associated with it
func DeleteSeason(ID string) {
	collection := database.DbClient.Database("test").Collection("seasons")
	filter := bson.M{
		"id": ID,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
	database.DbClient.Database("test").Collection("schedules").DeleteMany(context.TODO(), bson.M{
		"season": ID,
	})
}
