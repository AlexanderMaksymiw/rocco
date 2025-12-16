package main

import (
	"log"
	"time"
)

func StartReminderCron() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			runReminderCheck()
		}
	}()
}

func runReminderCheck() {
	now := time.Now().Unix()

	rows, err := db.Query(`
		SELECT id, user_id, title, message
		FROM reminder
		WHERE resolved = 0
		  AND notified = 0
		  AND due_date <= ?
	`, now)
	if err != nil {
		log.Println("cron: failed to query reminders:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID int
		var title, message string

		if err := rows.Scan(&id, &userID, &title, &message); err != nil {
			continue
		}

		// TODO: send notification (email, push, whatever)
		log.Println("NOTIFY:", userID, title)

		_, err := db.Exec(`
			UPDATE reminder
			SET notified = 1
			WHERE id = ?
		`, id)
		if err != nil {
			log.Println("cron: failed to mark notified:", err)
		}
	}
}
