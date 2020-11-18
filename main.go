package main

import (
	"fmt"
	"wilhelmiina/database"
	"wilhelmiina/schedule"
)

func main() {
	database.Init()
	subjects, _ := schedule.LoadSubjects()
	fmt.Println(subjects)
	database.Close() // Close database connection
}
