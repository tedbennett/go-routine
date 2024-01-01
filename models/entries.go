package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GoalEntry struct {
	Date   time.Time `json:"date"`
	GoalId uuid.UUID `json:"goal_id"`
}

func FetchEntries(db *sql.DB, goalIds []uuid.UUID) ([]GoalEntry, error) {
	sql := "SELECT date, goal_id FROM goal_entry WHERE goal_id IN (?" + strings.Repeat(",?", len(goalIds)-1) + ") ORDER BY date ASC"
	ids := Map(goalIds, func(id uuid.UUID) string { return id.String() })
	args := []interface{}{}
	for _, id := range ids {
		args = append(args, id)
	}
	rows, err := db.Query(sql, args...)
	// Exit if the SQL doesn't work for some reason

	if err != nil {
		return nil, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	entries := []GoalEntry{}
	for rows.Next() {
		entry := GoalEntry{}
		date := ""
		err = rows.Scan(&date, &entry.GoalId)
		// Exit if we get an error
		if err != nil {
			return nil, err
		}
		parsed, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, err
		}
		entry.Date = parsed
		entries = append(entries, entry)
	}
	return entries, nil
}

func InsertEntry(db *sql.DB, date time.Time, goalId uuid.UUID) (int64, error) {
	sql := "INSERT INTO goal_entry(goal_id, date) VALUES(?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()
	res, err2 := stmt.Exec(goalId.String(), date.Format(time.RFC3339))
	if err2 != nil {
		return 0, err2
	}
	return res.LastInsertId()
}

func DeleteEntry(db *sql.DB, goalId uuid.UUID, date time.Time) (int64, error) {
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
