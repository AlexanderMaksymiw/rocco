package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type MaintenanceRecordWithType struct {
	Record MaintenanceRecords
	Type   MaintenanceType
}

type MaintenanceType struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	IntervalMiles  int    `json:"interval_miles"`
	IntervalMonths int    `json:"interval_months"`
	Notes          string `json:"notes"`
}

type MaintenanceStatusResponse struct {
	TypeID         int    `json:"type_id"`
	Name           string `json:"name"`
	MilesRemaining int    `json:"miles_remaining"`
	TimeRemaining  int64  `json:"time_remaining"`
}

const secondsPerMonth = 2629746 // calculated by the average number of days per month 30.43

func CalculateMaintenanceStatus(
	carOdometer int,
	recMileage int,
	recDate int64,
	intervalMiles int,
	intervalMonths int,
) (int, int64) {
	milesSince := carOdometer - recMileage
	milesRemaining := intervalMiles - milesSince

	intervalSeconds := int64(intervalMonths) * secondsPerMonth
	dueDate := recDate + intervalSeconds
	now := time.Now().Unix()
	timeRemaining := dueDate - now

	return milesRemaining, timeRemaining

}

func GetMaintenanceStatus(db *sql.DB, userID int) ([]MaintenanceStatusResponse, error) {
	var car Car
	err := db.QueryRow(`SELECT id, odometer FROM car WHERE owner_id = ?`, userID).Scan(&car.ID, &car.Odometer)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
	 SELECT r.type_id, r.date_done, r.mileage, t.name, t.interval_miles, t.interval_months
	 FROM maintenance_records r
	 JOIN maintenance_types t ON r.type_id = t.id
	 WHERE r.car_id = ?
	 `, car.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MaintenanceStatusResponse

	for rows.Next() {
		var (
			typeID     int
			dateDone   int64
			recMileage int
			name       string
			intMiles   int
			intMonths  int
		)

		rows.Scan(&typeID, &dateDone, &recMileage, &name, &intMiles, &intMonths)
		fmt.Println("Type:", name, "intervalMiles:", intMiles, "recMileage:", recMileage)

		milesRem, timeRem := CalculateMaintenanceStatus(
			car.Odometer,
			recMileage,
			dateDone,
			intMiles,
			intMonths,
		)

		results = append(results, MaintenanceStatusResponse{
			TypeID:         typeID,
			Name:           name,
			MilesRemaining: milesRem,
			TimeRemaining:  timeRem,
		})
	}

	return results, nil
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
