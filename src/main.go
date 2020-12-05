package main

import (
	"wilhelmiina/api"
	"wilhelmiina/database"
)

func main() {
	database.Init() // Start database connection
	// user.CreateUser("admin", 999, "Aapo Harju", "test_admin_128738189@test.com", "admin")
	api.StartServer()
	database.Close() // Close database connection
}
