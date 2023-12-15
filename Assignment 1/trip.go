package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func tripmain() {
	r := mux.NewRouter()

	r.HandleFunc("/trips", publishTrip).Methods("POST")
	r.HandleFunc("/trips/{tripID}", getTripDetails).Methods("GET")
	r.HandleFunc("/trips/{tripID}/enroll", enrollInTrip).Methods("POST")
	r.HandleFunc("/trips/{tripID}", cancelTrip).Methods("DELETE")
	r.HandleFunc("/trips/user/{userID}", getPastTripsForUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func publishTrip(w http.ResponseWriter, r *http.Request) {
	var newTrip Trip

	// Decode the JSON request body into a Trip struct
	if err := json.NewDecoder(r.Body).Decode(&newTrip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assuming you've established a database connection
	db := ConnectDB()
	defer db.Close()

	// Insert the new trip details into the 'trips' table
	stmt, err := db.Prepare("INSERT INTO trips (CarOwnerID, PickupLocation, AltPickupLocation, StartTravelTime, Destination, AvailableSeats, EnrolledPassengers, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the trip details
	res, err := stmt.Exec(newTrip.CarOwnerID, newTrip.PickupLocation, newTrip.AltPickupLocation, newTrip.StartTravelTime, newTrip.Destination, newTrip.AvailableSeats, newTrip.EnrolledPassengers, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the last inserted trip
	lastID, _ := res.LastInsertId()
	newTrip.TripID = int(lastID)

	// Respond with the created trip in JSON format
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrip)
}

func getTripDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	db := ConnectDB()
	defer db.Close()

	var trip Trip
	err := db.QueryRow("SELECT TripID, CarOwnerID, PickupLocation, AltPickupLocation, StartTravelTime, Destination, AvailableSeats, EnrolledPassengers, CreatedAt FROM trips WHERE TripID = ?", tripID).
		Scan(&trip.TripID, &trip.CarOwnerID, &trip.PickupLocation, &trip.AltPickupLocation, &trip.StartTravelTime, &trip.Destination, &trip.AvailableSeats, &trip.EnrolledPassengers, &trip.CreatedAt)
	if err != nil {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func enrollInTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	var enrollment struct {
		UserID int `json:"userID"`
		// Other necessary fields for enrollment
	}
	err := json.NewDecoder(r.Body).Decode(&enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the trip details from the database
	var trip Trip
	err = db.QueryRow("SELECT * FROM trips WHERE TripID = ?", tripID).Scan(&trip.TripID, &trip.CarOwnerID, &trip.PickupLocation, &trip.AltPickupLocation, &trip.StartTravelTime, &trip.Destination, &trip.AvailableSeats, &trip.EnrolledPassengers, &trip.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if there are available seats
	if trip.AvailableSeats <= 0 {
		http.Error(w, "No available seats", http.StatusBadRequest)
		return
	}

	// Update the trip's EnrolledPassengers list with the new user ID
	trip.EnrolledPassengers = append(trip.EnrolledPassengers, enrollment.UserID)

	// Decrement available seats
	trip.AvailableSeats--

	// Update the trip details in the database
	stmt, err := db.Prepare("UPDATE trips SET AvailableSeats=?, EnrolledPassengers=? WHERE TripID=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(trip.AvailableSeats, trip.EnrolledPassengers, tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func cancelTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tripID := params["tripID"]

	stmt, err := db.Prepare("DELETE FROM trips WHERE TripID=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getPastTripsForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	// Query trips table to get past trips for a specific user
	rows, err := db.Query("SELECT TripID, CarOwnerID, PickupLocation, AltPickupLocation, StartTravelTime, Destination, AvailableSeats, EnrolledPassengers, CreatedAt FROM trips WHERE CarOwnerID = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pastTrips []Trip
	for rows.Next() {
		var trip Trip
		err := rows.Scan(&trip.TripID, &trip.CarOwnerID, &trip.PickupLocation, &trip.AltPickupLocation, &trip.StartTravelTime, &trip.Destination, &trip.AvailableSeats, &trip.EnrolledPassengers, &trip.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pastTrips = append(pastTrips, trip)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pastTrips)
}
