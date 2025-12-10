package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaintenanceRecords struct {
	ID       int64  `json:"id"`
	CarId    int64  `json:"car_id"`
	TypeId   int64  `json:"type_id"`
	DateDone string `json:"date_done"`
	Mileage  int64  `json:"mileage"`
	Notes    string `json:"notes"`
}

type MaintenanceRequest struct {
	TypeId   int64  `json:"type_id"`
	DateDone string `json:"date_done"`
	Mileage  int64  `json:"mileage"`
	Notes    string `json:"notes"`
}

type UpdateMaintenanceRequest struct {
	TypeID   *int64  `json:"type_id"`
	DateDone *string `json:"date_done"`
	Mileage  *int64  `json:"mileage"`
	Notes    *string `json:"notes"`
}

func getCarMaintenace(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	var carID int
	err := db.QueryRow("SELECT id FROM CAR WHERE owner_id = ?", userID).Scan(&carID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "car doesn't exist, add car to your account"})
		return
	}

	rows, err := db.Query(`
	SELECT id, car_id, type_id, date_done, mileage, notes
	FROM maintenance_records
	WHERE car_id = ?`, carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no car maintenance records found"})
		return
	}
	defer rows.Close()

	maintenanceList := []MaintenanceRecords{}

	for rows.Next() {
		var record MaintenanceRecords
		err := rows.Scan(
			&record.ID,
			&record.CarId,
			&record.TypeId,
			&record.DateDone,
			&record.Mileage,
			&record.Notes,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan maintenance records"})
			return
		}
		maintenanceList = append(maintenanceList, record)
	}

	if len(maintenanceList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"maintenance": []MaintenanceRecords{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"maintenance": maintenanceList})
}

func updateCarMaintenance(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	userID := userIDRaw.(int)

	maintenanceIDParam := c.Param("id")
	maintenanceID, err := strconv.ParseInt(maintenanceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse id parameter to int"})
		return
	}

	var request UpdateMaintenanceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var existing MaintenanceRecords
	err = db.QueryRow(`
	SELECT m.id, m.car_id, m.type_id, m.date_done, m.mileage, m.notes
	FROM maintenance_records m
	JOIN car c ON m.car_id = c.id
	Where m.id = ? AND c.owner_id = ? `,
		maintenanceID, userID).Scan(
		&existing.ID,
		&existing.CarId,
		&existing.TypeId,
		&existing.DateDone,
		&existing.Mileage,
		&existing.Notes,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no maintenance record found"})
		return
	}

	typeID := existing.TypeId
	if request.TypeID != nil {
		typeID = *request.TypeID
	}

	var dateDone int64
	if request.DateDone != nil {
		parsed, err := parseDateHelper(request.DateDone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_done"})
			return
		}
		dateDone = parsed
	} else {
		// Convert the existing string date to Unix if needed
		parsedExisting, err := parseDateHelper(&existing.DateDone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse existing date_done"})
			return
		}
		dateDone = parsedExisting
	}

	mileage := existing.Mileage
	if request.Mileage != nil {
		mileage = *request.Mileage
	}

	notes := existing.Notes
	if request.Notes != nil {
		notes = *request.Notes
	}

	_, err = db.Exec(`
	UPDATE maintenance_records SET
	type_id = ?,
	date_done = ?,
	mileage = ?,
	notes = ?
	Where id = ? AND car_id = ?`,
		typeID,
		dateDone,
		mileage,
		notes,
		userID,
		maintenanceID,
		existing.CarId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update maintenance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance record successfully updated"})
}

func addCarMaintenance(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID not found"})
		return
	}

	userID := userIDRaw.(int)

	var carID int64
	err := db.QueryRow("SELECT id FROM car WHERE owner_id = ?", userID).Scan(&carID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no car found for user"})
		return
	}

	var request MaintenanceRecords
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	parsed, err := parseDateHelper(&request.DateDone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format use (YYYY-MM-DD)"})
		return
	}

	dateDone := parsed

	recordInsert, err := db.Exec(`
	INSERT INTO maintenance_records
	(car_id, type_id, date_done, mileage, notes)
	VALUES (?, ?, ?, ?, ?)`,
		carID,
		request.TypeId,
		dateDone,
		request.Mileage,
		request.Notes)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create maintenance record"})
		return
	}

	ID, _ := recordInsert.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message": "maintenance record added successfully",
		"id":      ID,
	})

}
