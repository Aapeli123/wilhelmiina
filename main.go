package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/Aapeli123/wilhelmiina/user"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Aapeli123/wilhelmiina/api"
	"github.com/Aapeli123/wilhelmiina/database"
)

func main() {
	database.Init() // Start database connection
	admin, err := user.DoesAdminExist()
	if err != nil {
		panic(err)
	}
	if !admin { // Create an temporary admin account for users to access the service
		fmt.Println("IMPORTANT:\nThere are no admin users in the database.\nCreating a temporary admin user:")
		un, pw, err := credentials()
		if err != nil {
			panic(err)
		}
		_, err = user.CreateTemporaryAdmin(un, pw)
		if err != nil {
			panic(err)
		}
	}
	api.StartServer()
	database.Close() // Close database connection
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
