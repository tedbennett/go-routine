package models

import (
	"database/sql"
	"slices"
	"tb/goals/utils"

	"github.com/google/uuid"
)

// Models
type Goal struct {
	Id      uuid.UUID   `json:"id"`
	Entries []GoalEntry `json:"entries"`
	Title   string      `json:"title"`
	Color   string      `json:"color"`
}

func FetchGoal(db *sql.DB, userId uuid.UUID, goalId uuid.UUID) (Goal, error) {
	sql := "SELECT id, title, color FROM goal WHERE user_id = ? AND id = ?"
	goal := Goal{}
	err := db.QueryRow(sql, userId, goalId).Scan(&goal.Id, &goal.Title, &goal.Color)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return Goal{}, err
	}
	// make sure to cleanup when the program exits
	ids := []uuid.UUID{goal.Id}
	entries, _ := FetchEntries(db, ids)
	for _, entry := range entries {
		goal.Entries = append(goal.Entries, entry)
	}
	return goal, nil
}

func FetchGoalColor(db *sql.DB, goalId uuid.UUID) (string, error) {
	color := "#000000"
	sql := "SELECT color FROM goal WHERE id = ?"
	err := db.QueryRow(sql, goalId).Scan(&color)
	if err != nil {
		return "", err
	}
	return color, nil
}

func FetchGoals(db *sql.DB, userId uuid.UUID) ([]Goal, error) {
	sql := "SELECT id, title, color FROM goal WHERE user_id = ?"
	rows, err := db.Query(sql, userId)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return nil, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	goals := []Goal{}
	for rows.Next() {
		goal := Goal{}
		err2 := rows.Scan(&goal.Id, &goal.Title, &goal.Color)
		// Exit if we get an error
		if err2 != nil {
			return nil, err2
		}
		goals = append(goals, goal)
	}
	ids := utils.Map(goals, func(g Goal) uuid.UUID { return g.Id })
	entries, _ := FetchEntries(db, ids)
	for _, entry := range entries {
		idx := slices.IndexFunc(goals, func(g Goal) bool { return g.Id == entry.GoalId })
		if idx == -1 {
			continue
		}
		goals[idx].Entries = append(goals[idx].Entries, entry)
	}
	return goals, nil
}

func InsertGoal(db *sql.DB, userId uuid.UUID, title string, color string) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO goal(id, title, user_id, color) VALUES(?, ?, ?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return uuid.Nil, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id.String(), title, userId.String(), color)
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
