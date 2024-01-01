package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GoalEntry struct {
	Date   Date      `json:"date"`
	GoalId uuid.UUID `json:"goal_id"`
}

func (g *GoalEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date string `json:"date"`
	}{
		Date: g.Date.String(),
	})
}

func FetchEntries(db *sql.DB, goalIds []uuid.UUID) ([]GoalEntry, error) {
	sql := "SELECT date, goal_id FROM goal_entry WHERE goal_id IN (?" + strings.Repeat(",?", len(goalIds)-1) + ")"
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
		err2 := rows.Scan(&entry.Date, &entry.GoalId)
		// Exit if we get an error
		if err2 != nil {
			return nil, err2
		}
		entries = append(entries, entry)
	}
	fmt.Println(entries)
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
