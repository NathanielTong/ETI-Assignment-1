CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS carpooling;

-- Use the created database
USE carpooling;
-- Create the users table
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
);
CREATE TABLE IF NOT EXISTS trips (
    TripID INT AUTO_INCREMENT PRIMARY KEY,
    CarOwnerID INT,
    PickupLocation VARCHAR(255),
    AltPickupLocation VARCHAR(255),
    StartTravelTime DATETIME,
    Destination VARCHAR(255),
    AvailableSeats INT,
    EnrolledPassengers JSON,
    CreatedAt DATETIME
);