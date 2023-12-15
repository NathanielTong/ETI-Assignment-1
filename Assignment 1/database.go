package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	UserID        int       `json:"userID"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	MobileNumber  string    `json:"mobileNumber"`
	Email         string    `json:"email"`
	IsCarOwner    bool      `json:"isCarOwner"`
	DriverLicense string    `json:"driverLicense,omitempty"`
	CarPlateNum   string    `json:"carPlateNum,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Trip struct {
	TripID             int       `json:"tripID"`
	CarOwnerID         int       `json:"carOwnerID"`
	PickupLocation     string    `json:"pickupLocation"`
	AltPickupLocation  string    `json:"altPickupLocation"`
	StartTravelTime    time.Time `json:"startTravelTime"`
	Destination        string    `json:"destination"`
	AvailableSeats     int       `json:"availableSeats"`
	EnrolledPassengers []int     `json:"enrolledPassengers"`
	CreatedAt          time.Time `json:"createdAt"`
}

func ConnectDB() *sql.DB {
	connStr := "user:password@tcp(localhost:3306)/carpooling" // Update with your database connection details
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InitializeDatabase() {
	db = ConnectDB()
	defer db.Close()

	// Create tables if they don't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			UserID INT AUTO_INCREMENT PRIMARY KEY,
			FirstName VARCHAR(50) NOT NULL,
			LastName VARCHAR(50) NOT NULL,
			MobileNumber VARCHAR(15) NOT NULL,
			Email VARCHAR(100) NOT NULL,
			IsCarOwner BOOLEAN NOT NULL,
			DriverLicense VARCHAR(20),
			CarPlateNum VARCHAR(20),
			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func InitializeDatabasetrip() {
	db = ConnectDB()
	defer db.Close()

	// Create 'trips' table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS trips (
		TripID INT AUTO_INCREMENT PRIMARY KEY,
		CarOwnerID INT,
		PickupLocation VARCHAR(255),
		AltPickupLocation VARCHAR(255),
		StartTravelTime DATETIME,
		Destination VARCHAR(255),
		AvailableSeats INT,
		EnrolledPassengers JSON,
		CreatedAt DATETIME
	)`)
	if err != nil {
		log.Fatal(err)
	}
}
