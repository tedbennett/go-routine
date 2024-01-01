package handlers

import (
	"database/sql"
	"net/http"
	"tb/goals/models"

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
		entries, err := models.FetchEntries(db, id)
		if err != nil {
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
			return c.NoContent(http.StatusNotFound)
		}
		_, err = models.InsertEntry(db, id)
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
		now := models.DateNow()
		_, err2 := models.DeleteEntry(db, id, now)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}
