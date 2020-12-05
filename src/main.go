package main

import (
	"wilhelmiina/api"
	"wilhelmiina/database"
)

func main() {
	database.Init() // Start database conn
	api.StartServer()
	database.Close() // Close database connection
}
