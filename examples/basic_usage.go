package main

import (
	"fmt"
	"log"
	// "github.com/dreulavelle/GoDebridAPI/api"   // Be sure to import this!
)

func main() {
	// Initialize the client
	client := api.HttpClient(api.GetApiKey())

	// Fetch user details
	user, err := client.RdGetUser()
	if err != nil {
		log.Fatalf("Error fetching user details: %v", err)
	}
	fmt.Printf("User Details: %+v\n", user)

	// Add more examples for other API calls
}
