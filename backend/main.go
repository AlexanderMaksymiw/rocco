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
	router.DELETE("/delete-account", jwtRequired(), DeleteAccount)
	router.POST("/add-car", jwtRequired(), AddCar)
	router.GET("/car", jwtRequired(), getCar)
	router.PATCH("/update-car-info", jwtRequired(), updateCarInfo)
	router.DELETE("/delete-car", jwtRequired(), deleteCar)
	router.GET("/get-car-maintenace", jwtRequired(), getCarMaintenace)
	router.PATCH("/update-car-maintenace/:id", jwtRequired(), updateCarMaintenance)
	router.POST("/add-maintenance", jwtRequired(), addCarMaintenance)
	router.Run("localhost:8080")
}
