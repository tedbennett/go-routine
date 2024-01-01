package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"tb/goals/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		id, err := uuid.Parse(userId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		goals, err := models.FetchUser(db, id)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, goals)
	}
}

func PostUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		id, err := models.InsertUser(db, name)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, id)
	}
}

func PatchUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.Param("userId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		name := c.FormValue("name")
		_, err2 := models.UpdateUser(db, id, name)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, id)
	}
}
func DeleteUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")
		id, err := uuid.Parse(userId)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		_, err2 := models.DeleteUser(db, id)
		if err2 != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}
