package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// Models
type User struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FetchUser(db *sql.DB, id uuid.UUID) (User, error) {
	sql := "SELECT id, name from user WHERE id = ?"
	user := User{}
	err := db.QueryRow(sql, id).Scan(&user.Id, &user.Name)

	if err != nil {
		return User{}, err
	}
	return user, nil
}

func InsertUser(db *sql.DB, name string) (uuid.UUID, error) {
	id := uuid.New()
	sql := "INSERT INTO user(id, name) VALUES(?, ?)"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return uuid.Nil, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id.String(), name)
	if err2 != nil {
		return uuid.Nil, err2
	}
	_, err3 := res.LastInsertId()
	return id, err3
}

func UpdateUser(db *sql.DB, id uuid.UUID, name string) (int64, error) {
	sql := "UPDATE user SET name = ? WHERE id = ?"
	stmt, err := db.Prepare(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return 0, err
	}
	// make sure to cleanup when the program exits
	defer stmt.Close()

	res, err2 := stmt.Exec(id.String(), name)
	if err2 != nil {
		return 0, err2
	}
	rows, err3 := res.RowsAffected()
	return rows, err3
}

func DeleteUser(db *sql.DB, id uuid.UUID) (int64, error) {
	sql := "DELETE FROM user WHERE id = ?"
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
