package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{userID}", getUser).Methods("GET")
	r.HandleFunc("/users/{userID}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{userID}", deleteUser).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler to create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	// Decode the JSON request body into a User struct
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := ConnectDB()
	defer db.Close()

	// Insert the new user's details into the 'users' table
	_, err = db.Exec("INSERT INTO users (FirstName, LastName, MobileNumber, Email, IsCarOwner, DriverLicense, CarPlateNum, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		newUser.FirstName, newUser.LastName, newUser.MobileNumber, newUser.Email, newUser.IsCarOwner, newUser.DriverLicense, newUser.CarPlateNum, time.Now())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Respond with the created user profile
	json.NewEncoder(w).Encode(newUser)
}

// Handler to get user details
func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db := ConnectDB()
	defer db.Close()

	var user User
	err := db.QueryRow("SELECT UserID, FirstName, LastName, MobileNumber, Email, IsCarOwner, DriverLicense, CarPlateNum, CreatedAt FROM users WHERE UserID = ?", userID).
		Scan(&user.UserID, &user.FirstName, &user.LastName, &user.MobileNumber, &user.Email, &user.IsCarOwner, &user.DriverLicense, &user.CarPlateNum, &user.CreatedAt)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the fetched user profile in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handler to update user profile
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := ConnectDB()
	defer db.Close()

	// Update the user's details in the 'users' table
	_, err = db.Exec("UPDATE users SET FirstName=?, LastName=?, MobileNumber=?, Email=?, IsCarOwner=?, DriverLicense=?, CarPlateNum=? WHERE UserID=?",
		updatedUser.FirstName, updatedUser.LastName, updatedUser.MobileNumber, updatedUser.Email, updatedUser.IsCarOwner, updatedUser.DriverLicense, updatedUser.CarPlateNum, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with ID %s updated successfully", userID)
}

// Handler to delete user account
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db := ConnectDB()
	defer db.Close()

	// Delete the user from the 'users' table
	_, err := db.Exec("DELETE FROM users WHERE UserID=?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with ID %s deleted successfully", userID)
}
