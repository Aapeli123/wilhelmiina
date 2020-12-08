package schedule

import (
	"context"
	"fmt"
	"wilhelmiina/database"
	"wilhelmiina/user"

	"github.com/google/uuid"
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
	CourseID      string
}

// Group reprsents a group of students that has a subject
type Group struct {
	CourseID string
	Name     string
	Position int8
	Teacher  string
	Students []string
	GroupID  string
	SeasonID string
}

// AddStudent adds a new user to the group and saves change to database
func (g *Group) AddStudent(user user.User) {
	g.Students = append(g.Students, user.UUID)
	filter := bson.M{
		"groupid": g.GroupID,
	}
	collection := database.DbClient.Database("test").Collection("groups")
	collection.FindOneAndReplace(context.TODO(), filter, *g)
}

// AddGroup creates a new group and saves it to database
func AddGroup(c Course, teacher user.User, position int8, season Season) (Group, error) {
	prevGroups, err := GetGroupsForCourse(c.CourseID)
	if err != nil {
		return Group{}, err
	}
	amountOfGroups := len(prevGroups)
	groupName := fmt.Sprintf("%s.%d", c.NameShort, amountOfGroups+1)
	group := Group{
		Name:     groupName,
		GroupID:  uuid.New().String(),
		CourseID: c.CourseID,
		Position: position,
		Teacher:  teacher.UUID,
		Students: []string{},
		SeasonID: season.ID,
	}

	collection := database.DbClient.Database("test").Collection("groups")
	_, err = collection.InsertOne(context.TODO(), group)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

// GetGroupsForCourse gets all groups for specific course
func GetGroupsForCourse(CourseID string) ([]Group, error) {
	filter := bson.M{
		"courseid": CourseID,
	}
	cur, err := database.DbClient.Database("test").Collection("groups").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	groups := []Group{}
	for cur.Next(context.TODO()) {
		var elem Group
		cur.Decode(&elem)
		groups = append(groups, elem)
	}
	return groups, nil
}

// GetGroup gets group based on group id
func GetGroup(groupID string) (Group, error) {
	filter := bson.M{
		"groupid": groupID,
	}
	var group Group
	err := database.DbClient.Database("test").Collection("groups").FindOne(context.TODO(), filter).Decode(&group)
	if err != nil {
		return Group{}, nil
	}
	return group, nil
}

// AddCourse adds a new course and saves it to the database
func AddCourse(name string, desc string, number int, subject Subject) (Course, error) {
	course := Course{
		Name:          name,
		Description:   desc,
		CourseSubject: subject,
		Number:        number,
		NameShort:     fmt.Sprintf("%s%d", subject.Name, number),
		CourseID:      uuid.New().String(),
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
func GetCourse(id string) (Course, error) {
	collection := database.DbClient.Database("test").Collection("courses")
	filter := bson.M{
		"courseid": id,
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

// DeleteCourse removes a course from database
func DeleteCourse(courseID string) {
	collection := database.DbClient.Database("test").Collection("courses")
	filter := bson.M{
		"courseid": courseID,
	}
	collection.FindOneAndDelete(context.TODO(), filter)
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
