package main

import (
	"net/http"

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

type MaintenaceRequest struct {
	TypeId   int64  `json:"type_id"`
	DateDone string `json:"date_done"`
	Mileage  int64  `json:"mileage"`
	Notes    string `json:"notes"`
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
