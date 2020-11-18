package schedule

import (
	"context"
	"fmt"
	"wilhelmiina/database"
	"wilhelmiina/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Course represents a single course that the user can pick
type Course struct {
	Name          string
	Number        int
	NameShort     string
	Description   string
	CourseSubject Subject
}

// Group reprsents a group of students that has a subject
type Group struct {
	Course   Course
	Name     string
	Position int8
	Teacher  user.Teacher
	Students []user.User
}

// AddCourse adds a new course and saves it to the database
func AddCourse(name string, desc string, number int, subject Subject) (Course, error) {
	course := Course{
		Name:          name,
		Description:   desc,
		CourseSubject: subject,
		Number:        number,
		NameShort:     fmt.Sprintf("%s%d", subject.Name, number),
	}
	res, err := doesCourseExist(course.NameShort)
	if err != nil {
		return Course{}, err
	}
	if !res {
		collection := database.DbClient.Database("test").Collection("courses")
		collection.InsertOne(context.TODO(), course)
		return course, nil
	}
	return Course{}, ErrDocExistsAlready
}

// GetCourse gets course data from database
func GetCourse(shortName string) (Course, error) {
	collection := database.DbClient.Database("test").Collection("courses")
	filter := bson.M{
		"nameshort": shortName,
	}
	var course Course
	err := collection.FindOne(context.TODO(), filter).Decode(&course)
	if err != nil {
		return Course{}, err
	}
	return course, nil
}

// LoadCourses loads the course data from the database
func LoadCourses() ([]Course, error) {
	var courses []Course
	cur, err := database.DbClient.Database("test").Collection("courses").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem Course
		cur.Decode(&elem)
		courses = append(courses, elem)
	}
	return courses, nil
}

func doesCourseExist(shortName string) (bool, error) {
	_, err := GetCourse(shortName)
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}
