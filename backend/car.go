package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Car struct {
	ID                 int            `json:"id"`
	Make               string         `json:"make"`
	Model              string         `json:"model"`
	Engine             string         `json:"engine"`
	Year               string         `json:"year"`
	Odometer           int            `json:"odometer"`
	Vin                sql.NullString `json:"vin"`
	Mot_due            sql.NullInt64  `json:"mot_due"`
	Tax_due            sql.NullInt64  `json:"tax_due"`
	Insured_until      sql.NullInt64  `json:"insured_until"`
	Owner_id           int            `json:"owner_id"`
}

type UpdateCarRequest struct {
	Make               *string `json:"make"`
	Model              *string `json:"model"`
	Engine             *string `json:"engine"`
	Year               *string `json:"year"`
	Odometer           *int64  `json:"odometer"`
	Vin                *string `json:"vin"`
	Mot_due            *string `json:"mot_due"`
	Tax_due            *string `json:"tax_due"`
	Insured_until      *string `json:"insured_until"`
}

type CarRequest struct {
	Make   string `json:"make" binding:"required"`
	Model  string `json:"model" binding:"required"`
	Engine string `json:"engine"`
	Year   string `json:"year"`
	Vin    string `json:"vin"`
}

func parseDateHelper(field *string) (int64, error) {
	if field == nil {
		return 0, nil
	}

	layout := "2006-01-02"
	t, err := time.Parse(layout, *field)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func getCar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	var car Car
	row := db.QueryRow(`
		SELECT id, make, model, engine, year, odometer, vin, mot_due, tax_due, insured_until, owner_id
		FROM car
		WHERE owner_id = ?
		ORDER BY id LIMIT 1
	`, userID)

	err := row.Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Engine,
		&car.Year,
		&car.Odometer,
		&car.Vin,
		&car.Mot_due,
		&car.Tax_due,
		&car.Insured_until,
		&car.Owner_id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No car found for this user"})
		return
	}

	c.IndentedJSON(http.StatusOK, car)
}

func AddCar(c *gin.Context) {
	// Get the logged-in user's ID from JWT
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	userID := userIDRaw.(int)

	// Bind incoming JSON to CarRequest struct
	var req CarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, make and model are required"})
		return
	}

	// Insert new car into DB
	res, err := db.Exec(`
        INSERT INTO car (make, model, engine, year, vin, owner_id)
        VALUES (?, ?, ?, ?, ?, ?)
    `, req.Make, req.Model, req.Engine, req.Year, req.Vin, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create car"})
		return
	}

	carID, _ := res.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Car added successfully",
		"car_id":  carID,
	})
}

func updateCarInfo(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	userID := userIDRaw.(int)

	var request UpdateCarRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dateFields := map[string]*string{
		"mot_due":           request.Mot_due,
		"tax_due":           request.Tax_due,
		"insured_until":     request.Insured_until,
	}

	parsedDates := make(map[string]int64)

	for name, field := range dateFields {
		if field == nil {
			parsedDates[name] = 0
			continue
		}

		parsed, err := parseDateHelper(field)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s invalid format (expected YYYY-MM-DD)", name),
			})
			return
		}

		parsedDates[name] = parsed
	}

	_, err := db.Exec(`
		UPDATE car SET
			odometer = ?,
			mot_due = ?,
			tax_due = ?,
			insured_until = ?
		WHERE owner_id = ?`,
		request.Odometer,
		parsedDates["mot_due"],
		parsedDates["tax_due"],
		parsedDates["insured_until"],
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update car information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car successfully updated!"})
}

func deleteCar(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	userID := userIDRaw.(int)

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start transaction"})
		return
	}

	_, err = tx.Exec("DELETE FROM car WHERE owner_id = ?", userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete car"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "car successfully deleted"})

}
