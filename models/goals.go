package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// Models
type Goal struct {
	Id      uuid.UUID   `json:"id"`
	Entries []GoalEntry `json:"entries"`
	Title   string      `json:"title"`
}

func FetchGoals(db *sql.DB) ([]Goal, error) {
	sql := "SELECT id, title FROM goal"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return nil, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	result := []Goal{}
	for rows.Next() {
		goal := Goal{}
		err2 := rows.Scan(&goal.Id, &goal.Title)
		// Exit if we get an error
		if err2 != nil {
			return nil, err2
		}
		result = append(result, goal)
	}
	return result, nil
}

func InsertGoal(db *sql.DB, title string) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO goal(id, title) VALUES(?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return uuid.Nil, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id.String(), title)
	if err2 != nil {
		return uuid.Nil, err2
	}
	_, err3 := res.LastInsertId()
	return id, err3
}

func UpdateGoal(db *sql.DB, id uuid.UUID, title string) (int64, error) {
	sql := "UPDATE goal SET title = ? WHERE id = ?"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id.String(), title)
	if err2 != nil {
		return 0, err2
	}
	rows, err3 := res.RowsAffected()
	return rows, err3
}

func DeleteGoal(db *sql.DB, id uuid.UUID) (int64, error) {
	sql := "DELETE FROM goal WHERE id = ?"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id)
	if err2 != nil {
		return 0, err2
	}
	return res.RowsAffected()
}
