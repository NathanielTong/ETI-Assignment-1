package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func consolemain() {
	var firstName, lastName, mobileNumber, email, driverLicense, carPlateNum, pickupLocation, altPickupLocation, destination string
	var isCarOwner, ispublishingtrip bool
	var carOwnerID, availableSeats, userID, tripID int
	var startTravelTime time.Time
	// Collect user input
	fmt.Println("Enter First Name:")
	fmt.Scanln(&firstName)
	fmt.Println("Enter Last Name:")
	fmt.Scanln(&lastName)
	fmt.Println("Enter Mobile Number:")
	fmt.Scanln(&mobileNumber)
	fmt.Println("Enter Email:")
	fmt.Scanln(&email)
	fmt.Println("are you a car owner? (true/false):")
	fmt.Scanln(&isCarOwner)

	// If the user is a car owner, collect additional details
	if isCarOwner {
		fmt.Println("Enter Driver's License:")
		fmt.Scanln(&driverLicense)
		fmt.Println("Enter Car Plate Number:")
		fmt.Scanln(&carPlateNum)
		fmt.Println("are you a publishing a trip? (true/false):")
		fmt.Scanln(&ispublishingtrip)
	}
	if ispublishingtrip {
		// Collect input for trip details
		fmt.Println("Enter Pickup Location:")
		fmt.Scanln(&pickupLocation)
		fmt.Println("Enter Alternative Pickup Location:")
		fmt.Scanln(&altPickupLocation)
		fmt.Println("Enter Start Travel Time:")
		fmt.Scanln(&startTravelTime)
		fmt.Println("Enter Destination:")
		fmt.Scanln(&destination)
		fmt.Println("Enter Available Seats:")
		fmt.Scanln(&availableSeats)
	}

	db := ConnectDB()
	defer db.Close()

	err := db.QueryRow("SELECT UserID, FROM users WHERE FirstName = ? AND LastName = ?", firstName, lastName).
		Scan(&userID)
	if err != nil {
		return
	}

	err1 := db.QueryRow("SELECT TripID FROM trips WHERE CarOwnerID = ?", userID).
		Scan(&tripID)
	if err1 != nil {
		return
	}
	createUserConsole(firstName, lastName, mobileNumber, email, isCarOwner, driverLicense, carPlateNum)
	getUserConsole(userID)
	updateUserConsoleinput(userID)
	deleteUserConsole(userID)
	publishTripConsole(carOwnerID, pickupLocation, altPickupLocation, startTravelTime, destination, availableSeats)
	getTripDetailsConsole(tripID)
	enrollInTripConsole(userID, tripID)
	cancelTripConsole(tripID)
}

func createUserConsole(firstName, lastName, mobileNumber, email string, isCarOwner bool, driverLicense, carPlateNum string) {
	user := User{
		FirstName:     firstName,
		LastName:      lastName,
		MobileNumber:  mobileNumber,
		Email:         email,
		IsCarOwner:    isCarOwner,
		DriverLicense: driverLicense,
		CarPlateNum:   carPlateNum,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8080/users" // Replace with your API endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("User created successfully")
	} else {
		fmt.Println("Failed to create user")
	}
}

func getUserConsole(userID int) {
	url := fmt.Sprintf("http://localhost:8080/users/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}

	fmt.Println("User Details:")
	fmt.Printf("UserID: %d\nName: %s %s\nMobile: %s\nEmail: %s\nIsCarOwner: %t\n", user.UserID, user.FirstName, user.LastName, user.MobileNumber, user.Email, user.IsCarOwner)
}

func updateUserConsoleinput(userID int) {
	var firstName, lastName, mobileNumber, email, driverLicense, carPlateNum string
	var isCarOwner bool
	updatedUser := User{
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		MobileNumber:  mobileNumber,
		Email:         email,
		IsCarOwner:    isCarOwner,
		DriverLicense: driverLicense,
		CarPlateNum:   carPlateNum,
		CreatedAt:     time.Now(),
	}
	// Collect updated user information
	fmt.Println("Enter First Name:")
	fmt.Scanln(&firstName)
	fmt.Println("Enter Last Name:")
	fmt.Scanln(&lastName)
	fmt.Println("Enter Mobile Number:")
	fmt.Scanln(&mobileNumber)
	fmt.Println("Enter Email:")
	fmt.Scanln(&email)
	fmt.Println("Are you a car owner? (true/false):")
	fmt.Scanln(&isCarOwner)

	// If the user is a car owner, collect additional details
	if isCarOwner {
		fmt.Println("Enter Driver's License:")
		fmt.Scanln(&driverLicense)
		fmt.Println("Enter Car Plate Number:")
		fmt.Scanln(&carPlateNum)
	}

	// Call the updateUser function with the collected information
	updateUserConsole(userID, updatedUser)
}

func updateUserConsole(userID int, updatedUser User) {
	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:8080/users/%d", userID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("User profile updated successfully")
	} else {
		fmt.Println("Failed to update user profile")
	}
}

func deleteUserConsole(userID int) {
	url := fmt.Sprintf("http://localhost:8080/users/%d", userID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("User account deleted successfully")
	} else {
		fmt.Println("Failed to delete user account")
	}
}

func publishTripConsole(carOwnerID int, pickupLocation, altPickupLocation string, startTravelTime time.Time, destination string, availableSeats int) {
	newTrip := Trip{
		CarOwnerID:        carOwnerID,
		PickupLocation:    pickupLocation,
		AltPickupLocation: altPickupLocation,
		StartTravelTime:   startTravelTime,
		Destination:       destination,
		AvailableSeats:    availableSeats,
	}

	jsonData, err := json.Marshal(newTrip)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8080/trips" // Replace with your API endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Trip published successfully")
	} else {
		fmt.Println("Failed to publish trip")
	}
}

func getTripDetailsConsole(tripID int) {
	url := fmt.Sprintf("http://localhost:8080/trips/%d", tripID) // Replace with your API endpoint
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var trip Trip
	if err := json.NewDecoder(resp.Body).Decode(&trip); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Trip Details:")
	fmt.Printf("TripID: %d\nCarOwnerID: %d\nPickup Location: %s\nAlt Pickup Location: %s\nStart Travel Time: %s\nDestination: %s\nAvailable Seats: %d\n",
		trip.TripID, trip.CarOwnerID, trip.PickupLocation, trip.AltPickupLocation, trip.StartTravelTime, trip.Destination, trip.AvailableSeats)
}

func enrollInTripConsole(tripID, userID int) {
	enrollment := struct {
		UserID int `json:"userID"`
	}{
		UserID: userID,
	}

	jsonData, err := json.Marshal(enrollment)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:8080/trips/%d/enroll", tripID) // Replace with your API endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Enrolled in the trip successfully")
	} else {
		fmt.Println("Failed to enroll in the trip")
	}
}

func cancelTripConsole(tripID int) {
	url := fmt.Sprintf("http://localhost:8081/trips/%d", tripID) // Replace with your API endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Trip canceled successfully")
	} else {
		fmt.Println("Failed to cancel trip")
	}
}

func getPastTripsForUserConsole(userID int) {
	url := fmt.Sprintf("http://localhost:8080/trips/user/%d", userID) // Replace with your API endpoint
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var pastTrips []Trip
	if err := json.NewDecoder(resp.Body).Decode(&pastTrips); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Past Trips:")
	for _, trip := range pastTrips {
		fmt.Printf("TripID: %d\nCarOwnerID: %d\nPickupLocation: %s\nDestination: %s\n\n", trip.TripID, trip.CarOwnerID, trip.PickupLocation, trip.Destination)
	}
}
