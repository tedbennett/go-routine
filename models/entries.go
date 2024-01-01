package models

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type GoalEntry struct {
	Date Date `json:"date"`
}

func (g *GoalEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date string `json:"date"`
	}{
		Date: g.Date.String(),
	})
}

func FetchEntries(db *sql.DB, goalId uuid.UUID) ([]GoalEntry, error) {
	sql := "SELECT date FROM goal_entry"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return nil, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	entries := []GoalEntry{}
	for rows.Next() {
		entry := GoalEntry{}
		err2 := rows.Scan(&entry.Date)
		// Exit if we get an error
		if err2 != nil {
			return nil, err2
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func InsertEntry(db *sql.DB, goalId uuid.UUID) (int64, error) {
	sql := "INSERT INTO goal_entry(goal_id, date) VALUES(?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()
	date := DateNow()
	res, err2 := stmt.Exec(goalId.String(), date.String())
	if err2 != nil {
		return 0, err2
	}
	return res.LastInsertId()
}

// Can only delete the most recent entry for now
func DeleteEntry(db *sql.DB, goalId uuid.UUID, date Date) (int64, error) {
	sql := "DELETE FROM goal_entry WHERE goal_id = ? AND date = ?"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()
	res, err2 := stmt.Exec(goalId.String(), date.String())
	if err2 != nil {
		return 0, err2
	}
	return res.LastInsertId()
}
