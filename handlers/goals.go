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
		USER, _ := uuid.Parse("5c1c0569-dcd1-4c0d-87f0-0d0c1debdd5b")
		goals, err := models.FetchGoals(db, USER)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, goals)
	}
}

func GetGoal(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		USER, _ := uuid.Parse("5c1c0569-dcd1-4c0d-87f0-0d0c1debdd5b")

		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		goals, err := models.FetchGoal(db, USER, id)
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

		USER, _ := uuid.Parse("5c1c0569-dcd1-4c0d-87f0-0d0c1debdd5b")
		id, err := models.InsertGoal(db, USER, title, "#000000")
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
