package user

import (
	"testing"
	"wilhelmiina/database"
)

// TestCreateUser Tests user creation
func TestCreateUser(t *testing.T) {
	database.Init()
	t.Log("Creating user...")
	user, err := CreateUser("test", 1, "test mac testface", "test@test.com", "test")
	if err != nil {
		t.Error(err)
	}
	r, _ := user.CheckPassword("test")
	r2, _ := user.CheckPassword("not test")
	if !r || r2 {
		t.Error("Password check failed")
	}
	DeleteUser(user.UUID)
	database.Close()
}

func TestGetUser(t *testing.T) {
	database.Init()
	t.Log("Creating user...")
	user, _ := CreateUser("test", 1, "test mac testface", "test@test.com", "test")
	alsoUser, err := GetUser(user.UUID)
	if err != nil {
		t.Error(err)
	}
	if user.UUID != alsoUser.UUID {
		t.Error("User ids don't match")
	}
	DeleteUser(user.UUID)
	database.Close()
}

func TestUserDelete(t *testing.T) {
	database.Init()
	user, _ := CreateUser("testuser2", 1, "test test", "test@test.com", "test")
	DeleteUser(user.UUID)
	r, _ := doesUserExist(user.UUID)
	if r {
		t.Error("User was not deleted")
	}
	database.Close()
}

func TestUserUpdates(t *testing.T) {
	database.Init()
	user, _ := CreateUser("testuser3", 1, "test test", "test@test.com", "test")
	UpdateEmail(user.UUID, "test2@test.com")
	user, _ = GetUser(user.UUID)
	if user.Email != "test2@test.com" {
		t.Error("User email was not changed")
	}
	UpdateUsername(user.UUID, "Yeet")
	user, _ = GetUser(user.UUID)
	if user.Username != "Yeet" {
		t.Error("Username was not updated")
	}
	UpdatePassword(user.UUID, "epicGamer")
	user, _ = GetUser(user.UUID)
	r1, _ := user.CheckPassword("test")
	r2, _ := user.CheckPassword("epicGamer")

	if !r2 || r1 {
		t.Error("Password was not changed")
	}
	DeleteUser(user.UUID)
	database.Close()
}

func TestGetUserByName(t *testing.T) {
	database.Init()
	t.Log("Creating user...")
	user, _ := CreateUser("test", 1, "test mac testface", "test@test.com", "test")
	alsoUser, err := GetUserByName(user.Username)
	if err != nil {
		t.Error(err)
	}
	if user.UUID != alsoUser.UUID {
		t.Error("User ids don't match")
	}
	DeleteUser(user.UUID)
	database.Close()
}
