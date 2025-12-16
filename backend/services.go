package main

import (
	"database/sql"
	"fmt"
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
