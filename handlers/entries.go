package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"tb/goals/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetEntries(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		entries, err := models.FetchEntries(db, []uuid.UUID{id})
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, entries)
	}
}

func PutEntry(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {

			fmt.Println(err)
			return c.NoContent(http.StatusNotFound)
		}
		date := c.FormValue("date")
		parsed, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}
		parsed = models.TruncateToDay(parsed)
		_, err = models.InsertEntry(db, parsed, id)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func DeleteEntry(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		now := models.TruncateToDay(time.Now())
		_, err2 := models.DeleteEntry(db, id, now)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}
