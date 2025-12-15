package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getReminders(c *gin.Context) {
	userID := c.GetInt("userID")

	var carID int
	err := db.QueryRow("SELECT id FROM car WHERE owner_id = ?", userID).Scan(&carID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "car not found"})
		return
	}

	rows, err := db.Query(`
	SELECT id, title, message, due_date, threshold_miles
	FROM reminder
	WHERE user_id = ? AND car_id = ? AND resolved = 0
	`, userID, carID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load reminders"})
	}
	defer rows.Close()

	reminders := []map[string]interface{}{}
	for rows.Next() {
		var id, thresholdMiles int
		var title, message string
		var dueDate int64

		rows.Scan(&id, &title, &message, &dueDate, &thresholdMiles)
		reminders = append(reminders, map[string]interface{}{
			"id":              id,
			"title":           title,
			"message":         message,
			"due_date":        dueDate,
			"threshold_miles": thresholdMiles,
		})
	}
	c.JSON(http.StatusOK, reminders)
}
