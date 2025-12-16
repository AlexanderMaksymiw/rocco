package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func CreateReminders(userID int, carID int) error {
	var car Car
	err := db.QueryRow(`
	SELECT id,
	odometer,
	mot_due,
	tax_due,
	insured_until
	FROM car WHERE owner_id = ?`,
		userID).Scan(&car.ID, &car.Odometer, &car.Mot_due, &car.Tax_due, &car.Insured_until)
	if err != nil {
		return err
	}

	rows, err := db.Query(`
	SELECT r.id, r.car_id, r.type_id, r.date_done, r.mileage,
				 t.name, t.interval_miles, t.interval_months
	FROM maintenance_records r
	JOIN maintenance_types t ON r.type_id = t.id
	WHERE r.car_id = ?
		AND r.date_done = (
				SELECT MAX(r2.date_done)
				FROM maintenance_records r2
				WHERE r2.car_id = r.car_id
					AND r2.type_id = r.type_id
		)
	`, carID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var recordsWithTypes []MaintenanceRecordWithType

	for rows.Next() {
		var rec MaintenanceRecords
		var mt MaintenanceType

		if err := rows.Scan(
			&rec.ID, &rec.CarId, &rec.TypeId, &rec.DateDone, &rec.Mileage,
			&mt.Name, &mt.IntervalMiles, &mt.IntervalMonths,
		); err != nil {
			return err
		}

		dateDoneInt, err := strconv.ParseInt(rec.DateDone, 10, 64)
		if err != nil {
			return err
		}

		milesInt := int(rec.Mileage)

		recordsWithTypes = append(recordsWithTypes, MaintenanceRecordWithType{
			Record: rec,
			Type:   mt,
		})

		milesRem, timeRem := CalculateMaintenanceStatus(
			car.Odometer,
			milesInt,
			dateDoneInt,
			mt.IntervalMiles,
			mt.IntervalMonths,
		)

		fmt.Println("Reminder for", mt.Name, "miles left:", milesRem, "time left:", timeRem)

	}

	milesThresholds := []int{2500, 1000, 500, 0}
	timeThresholds := []int64{
		30 * 86400, // 1 month
		14 * 86400, // 2 weeks
		7 * 86400,  // 1 week
		3 * 86400,  // 3 days
		0,          // due now
	}

	for _, rec := range recordsWithTypes {
		dateDoneInt, _ := strconv.ParseInt(rec.Record.DateDone, 10, 64)

		milesInt := int(rec.Record.Mileage)

		milesRem, timeRem := CalculateMaintenanceStatus(
			car.Odometer,
			milesInt,
			dateDoneInt,
			rec.Type.IntervalMiles,
			rec.Type.IntervalMonths,
		)

		var triggerMiles int = -1
		for _, t := range milesThresholds {
			if milesRem <= t {
				triggerMiles = t
				break
			}
		}

		var triggerTime int64 = -1
		for _, t := range timeThresholds {
			if timeRem <= t {
				triggerTime = t
				break
			}
		}

		if triggerMiles >= 0 || triggerTime >= 0 {

			var existingID int
			err := db.QueryRow(`
        SELECT id
        FROM reminder
        WHERE user_id = ?
          AND car_id = ?
          AND maintenance_type_id = ?
          AND resolved = 0
        LIMIT 1
    `, userID, carID, rec.Record.TypeId).Scan(&existingID)

			if err == nil {
				// reminder already exists, skip inserting
				continue
			}

			if err != sql.ErrNoRows {
				return err
			}

			_, err = db.Exec(`
    UPDATE reminder
    SET resolved = 1
    WHERE user_id = ?
      AND car_id = ?
      AND maintenance_type_id = ?
      AND resolved = 0
`, userID, carID, rec.Record.TypeId)

			if err != nil {
				return err
			}

			_, err = db.Exec(`
			INSERT INTO reminder
			(user_id, car_id, reminder_type, title, message, due_date, maintenance_record_id, maintenance_type_id, threshold_miles, resolved, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?)
			`,
				userID,
				carID,
				"maintenance",
				fmt.Sprintf("%s due soon", rec.Type.Name),
				fmt.Sprintf("Miles remaining: %d, Time remaining: %d seconds", milesRem, timeRem),
				time.Now().Unix()+timeRem,
				rec.Record.ID,
				rec.Record.TypeId,
				milesRem,
				time.Now().Unix(),
			)

			if err != nil {
				return err
			}
		}
	}
	return nil
}
