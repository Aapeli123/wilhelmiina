package schedule

import (
	"testing"

	"github.com/Aapeli123/wilhelmiina/database"
	"github.com/Aapeli123/wilhelmiina/user"

	"go.mongodb.org/mongo-driver/mongo"
)

// TESTS FOR SUBJECTS

func TestSubjectCreation(t *testing.T) {
	database.Init()
	s, err := AddSubject("TEST1", "Testiaine 1")
	if err != nil {
		t.Error(err)
	}
	if s.Name != "TEST1" || s.Description != "Testiaine 1" {
		t.Error("returned subject did not have the right name or description")
	}
	_, err = AddSubject("TEST1", "Testiaine also")
	if err != ErrDocExistsAlready {
		t.Error("Added subject even if there already was one with same name")
	}
	DeleteSubject("TEST1")
	database.Close()
}

func TestSubjectDeletion(t *testing.T) {
	database.Init()
	s, err := AddSubject("TEST2", "Testing")
	if err != nil {
		t.Error(err)
	}

	r, err := doesSubjectExist(s.Name)
	if err != nil {
		t.Error(err)
	}
	if !r {
		t.Error("Subject did not exist even though it should have")
	}

	DeleteSubject(s.Name)
	r, err = doesSubjectExist(s.Name)

	if err != nil {
		t.Error(err)
	}
	if r {
		t.Error("Subject was not deleted")
	}
	database.Close()
}

func TestGetSubject(t *testing.T) {
	database.Init()
	s, err := AddSubject("TEST2", "Testing")
	if err != nil {
		t.Error(err)
	}
	s, err = GetSubject(s.Name)
	if err != nil {
		t.Error(err)
	}
	if s.Description != "Testing" {
		t.Error("Returned subject description was wrong")
	}
	DeleteSubject("TEST2")
	s2, err := GetSubject("TEST2")
	if err != mongo.ErrNoDocuments {
		t.Error("GetSubject did not error when it should have. Error message: ", err, ". Returned subject: ", s2)
	}
	database.Close()
}

// TESTS FOR PROGRAM FLOW

func TestProgramFlow(t *testing.T) {
	database.Init()

	// create all necessary variables to the database
	subject, err := AddSubject("TEST", "Testiaine")
	if err != nil {
		t.Error("AddSubject failed:", err)
	}
	course, err := AddCourse("Testikurssi", "Testing testing", 1, subject)
	if err != nil {
		t.Error("AddCourse failed:", err)
	}

	teacher, err := user.CreateUser("testTeacher", 2, "Test Ope", "opettaja@test.com", "password")
	if err != nil {
		t.Error(err)
	}

	student, err := user.CreateUser("testStudent", 2, "Test Oppilas", "oppilas@test.com", "password")
	if err != nil {
		t.Error(err)
	}

	season, err := AddSeason("Testijakso", 0, 1)
	if err != nil {
		t.Error("AddSeason failed:", err)
	}

	teacherSchedule, err := AddSchedule(teacher.UUID, season.ID)
	if err != nil {
		t.Error("AddSchedule failed:", err)
	}

	studentSchedule, err := AddSchedule(student.UUID, season.ID)
	if err != nil {
		t.Error("AddSchedule failed:", err)
	}

	group, err := AddGroup(course, teacher, 1, season)
	if err != nil {
		t.Error("AddGroup failed:", err)
	}
	teacherSchedule.AddGroup(group.GroupID)

	group.AddStudent(student)
	studentSchedule.AddGroup(group.GroupID)

	// Test by reading the data from database
	groups, err := GetGroupsForCourse(course.CourseID)
	if err != nil {
		t.Error("GetGroupsForCourse failed", err)
	}
	if !(len(groups) > 0) {
		t.Error("There were no groups, there should have been, since one was just added")
	}

	// Check if our group was one of the returned
	found := false
	foundGroup := Group{}
	for _, element := range groups {
		if element.GroupID == group.GroupID {
			found = true
			foundGroup = element
			break
		}
	}

	if !found {
		t.Error("Testing group was not part of the groups returned")
	}

	dbCourse, err := GetCourse(foundGroup.CourseID)
	if err != nil {
		t.Log(foundGroup)
		t.Error("GetCourse failed", err)
	}
	if dbCourse.CourseID != course.CourseID {
		t.Error("GetCourse returned wrong course")
	}

	if len(foundGroup.Students) < 1 {
		t.Error("Found group did not have any students")
	}

	if foundGroup.Students[0] != student.UUID {
		t.Error("Database group had wrong student")
	}
	groupStudent, _ := user.GetUser(foundGroup.Students[0])
	if groupStudent.UUID != student.UUID {
		t.Error("Wrong uuid")
	}
	userSchedule, err := GetScheduleForUser(groupStudent.UUID, season.ID)
	if err != nil {
		t.Error("GetSchedule failed:", err)
	}
	if userSchedule.ScheduleID != studentSchedule.ScheduleID {
		t.Error("Get schedule returned wrong schedule")
	}

	schedulesSeason := userSchedule.Season
	dbSeason, err := GetSeason(schedulesSeason)
	if dbSeason.Name != "Testijakso" {
		t.Error("GetSeason returned the wrong season")
	}

	if userSchedule.Groups[0] != group.GroupID {
		t.Error("User schedule had the wrong group")
	}

	// Cleanup
	DeleteSubject(subject.Name)
	DeleteCourse(course.CourseID)
	DeleteGroup(group.GroupID)
	DeleteSeason(season.ID)
	_, err = GetSchedule(teacherSchedule.ScheduleID)

	user.DeleteUser(teacher.UUID)
	user.DeleteUser(student.UUID)

	database.Close()
}
