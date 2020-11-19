package user

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"wilhelmiina/database"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// User represents an user of the service
type User struct {
	Username     string
	Fullname     string
	Email        string
	PasswordHash string
	UUID         string
	LastLogin    int64
	Online       bool
	ScheduleIDs  []string
}

// Teacher represents a teacher
type Teacher struct {
	User
	ShortenedName string
}

var passwordConfig = hashParameters{
	keyLen:  128,
	memory:  64 * 1024,
	threads: 4,
	time:    3,
}

// CheckPassword Compares users password with another password
func (u *User) CheckPassword(password string) (bool, error) {
	parts := strings.Split(u.PasswordHash, "$")

	c := &hashParameters{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.keyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}

// AddSchedule adds the schedule id to user schedule ids list, then it updates database
func (u *User) AddSchedule(scheduleID string) {
	u.ScheduleIDs = append(u.ScheduleIDs, scheduleID)
	collection := database.DbClient.Database("test").Collection("users")
	filter := bson.M{
		"uuid": u.UUID,
	}
	collection.FindOneAndReplace(context.TODO(), filter, *u)
}

// CreateUser makes a new user and saves it to database.
// It returns the user.
func CreateUser(username string, fullName string, email string, password string) (User, error) {
	hashed, err := hashPassword(password, &passwordConfig)
	if err != nil {
		return User{}, err
	}
	user := User{
		UUID:         uuid.New().String(),
		Username:     username,
		Fullname:     fullName,
		Email:        email,
		PasswordHash: hashed,
		Online:       false,
		LastLogin:    time.Now().Unix(),
	}
	collection := database.DbClient.Database("test").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUser gets user based on user id
func GetUser(id string) (User, error) {
	collection := database.DbClient.Database("test").Collection("users")
	filter := bson.M{
		"uuid": id,
	}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// DeleteUser removes user from database
func DeleteUser(id string) {
	collection := database.DbClient.Database("test").Collection("users")
	collection.FindOneAndDelete(context.TODO(), bson.M{"uuid": id})
}

/** Functions that update users in some way: */

// UpdateUsername replaces users username on database with new username.
func UpdateUsername(id string, newUn string) {
	replaceUserData(id, "username", newUn)
}

// UpdateEmail replaces specific users email in database with the new email.
func UpdateEmail(id string, newEmail string) {
	replaceUserData(id, "email", newEmail)
}

// UpdatePassword updates the users password hash with the new passwords hash
func UpdatePassword(id string, newPass string) bool {
	hashed, err := hashPassword(newPass, &passwordConfig)
	if err != nil {
		return false
	}
	replaceUserData(id, "passwordhash", hashed)
	return true
}

// UpdateRealName updates users real name to database
func UpdateRealName(id string, newName string) {
	replaceUserData(id, "fullname", newName)
}

// GetUserByName returns an user based on username
func GetUserByName(name string) (User, error) {
	collection := database.DbClient.Database("test").Collection("users")
	filter := bson.M{
		"username": name,
	}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil

}

type hashParameters struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func hashPassword(password string, param *hashParameters) (string, error) {
	salt := make([]byte, 64)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, param.time, param.memory, param.threads, param.keyLen)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, param.memory, param.time, param.threads, b64salt, b64hash)
	return full, nil
}

// Replaces specific value from database for user specified by id
func replaceUserData(id string, propery string, newVal interface{}) {
	collection := database.DbClient.Database("test").Collection("users")
	filter := bson.M{
		"uuid": id,
	}
	collection.FindOneAndUpdate(context.TODO(), filter, bson.D{
		{"$set", bson.D{{propery, newVal}}},
	})
}

func doesUserExist(id string) (bool, error) {
	_, err := GetUser(id)
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}
