package schedule

import (
	"context"
	"wilhelmiina/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Subject represents a single subject
type Subject struct {
	Name        string
	Description string
}

// AddSubject adds a new subject to the database of subjects
func AddSubject(name string, desc string) (Subject, error) {
	collection := database.DbClient.Database("test").Collection("subjects")
	res, err := doesSubjectExist(name)
	if err != nil {
		return Subject{}, err
	}
	if !res {
		subject := Subject{
			Name:        name,
			Description: desc,
		}
		collection.InsertOne(context.TODO(), subject)
		return subject, nil
	}
	return Subject{}, ErrDocExistsAlready
}

// DeleteSubject removes a subject from database
func DeleteSubject(name string) {
	collection := database.DbClient.Database("test").Collection("subjects")
	filter := bson.M{
		"name": name,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
}

// GetSubject gets data for specific subject in the database
func GetSubject(name string) (Subject, error) {
	collection := database.DbClient.Database("test").Collection("subjects")
	filter := bson.M{
		"name": name,
	}
	var subject Subject
	err := collection.FindOne(context.TODO(), filter).Decode(&subject)
	if err != nil {
		return Subject{}, err
	}
	return subject, nil
}

func doesSubjectExist(name string) (bool, error) {
	_, err := GetSubject(name)
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}

// LoadSubjects Gets all subjects from database
func LoadSubjects() ([]Subject, error) {
	var subjects []Subject
	cur, err := database.DbClient.Database("test").Collection("subjects").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem Subject
		cur.Decode(&elem)
		subjects = append(subjects, elem)
	}
	return subjects, nil
}
