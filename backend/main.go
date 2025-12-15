package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	godotenv.Load()
}

func main() {
	var err error
	db, err = sql.Open("sqlite", `C:/Users/Alex/Documents/Rock-My-Rocco/rocco-database/scirocco.db`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		panic(err)
	}

	router := gin.Default()
	router.POST("/login", Login)
	router.POST("/signup", Signup)
	router.DELETE("/account", jwtRequired(), DeleteAccount)
	router.POST("/car", jwtRequired(), AddCar)
	router.GET("/car", jwtRequired(), getCar)
	router.PATCH("/car", jwtRequired(), updateCarInfo)
	router.DELETE("/car", jwtRequired(), deleteCar)
	router.GET("/maintenance", jwtRequired(), getCarMaintenance)
	router.POST("/maintenance", jwtRequired(), addCarMaintenance)
	router.PATCH("/maintenace/:id", jwtRequired(), updateCarMaintenance)
	router.GET("/maintenance/stats", jwtRequired(), getCarMaintenanceStats)
	router.POST("/maintenance/reminder", jwtRequired(), getCarMaintenanceReminders)
	router.Run("localhost:8080")
}
