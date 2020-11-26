package main

import (
	"fmt"
	"wilhelmiina/api"
	"wilhelmiina/database"
	"wilhelmiina/schedule"
)

func main() {
	database.Init()
	subjects, _ := schedule.LoadSubjects()
	api.StartServer()
	fmt.Println(subjects)
	database.Close() // Close database connection
}
