package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"tb/goals/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetGoals(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goals, err := models.FetchGoals(db)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, goals)
	}
}

func PostGoal(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		title := c.FormValue("title")
		id, err := models.InsertGoal(db, title)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, id)
	}
}

func PatchGoal(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		title := c.FormValue("title")
		_, err2 := models.UpdateGoal(db, id, title)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, id)
	}
}
func DeleteGoal(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		_, err2 := models.DeleteGoal(db, id)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}
