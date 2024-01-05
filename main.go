package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	database "tb/goals/db"
	"tb/goals/handlers"
	"tb/goals/models"
	"tb/goals/utils"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Templates

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplateRenderer(e *echo.Echo, paths ...string) {
	tmpl := &template.Template{}
	for i := range paths {
		template.Must(tmpl.ParseGlob(paths[i]))
	}
	t := newTemplate(tmpl)
	e.Renderer = t
}

func newTemplate(templates *template.Template) echo.Renderer {
	return &Template{
		Templates: templates,
	}
}

type IndexTemplateData struct {
	Goals []Goal
}

type Goal struct {
	Title   string
	Id      uuid.UUID
	Entries []Entry
}

type Entry struct {
	Complete bool
	Date     string
	Color    string
}

func mapGoals(goals []models.Goal) []Goal {
	return utils.Map(goals, func(g models.Goal) Goal {
		return Goal{Title: g.Title, Id: g.Id, Entries: buildEntries(g.Entries, g.Color)}
	})
}

func buildEntries(from []models.GoalEntry, color string) []Entry {
	entries := []Entry{}
	end := utils.TruncateDate(time.Now())
	start := end.AddDate(0, 0, -119)
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		complete := false
		for _, entry := range from {
			if utils.DateEqual(entry.Date, d) {
				complete = true
				break
			}
		}
		entries = append(entries, Entry{Date: d.Format(time.RFC3339), Complete: complete, Color: color})
	}
	return entries
}

func handleRoutine(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		title := c.FormValue("title")
		color := c.FormValue("color")
		USER, _ := uuid.Parse("5c1c0569-dcd1-4c0d-87f0-0d0c1debdd5b")
		id, err := models.InsertGoal(db, USER, title, color)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}
		data := Goal{Title: title, Id: id, Entries: buildEntries([]models.GoalEntry{}, color)}
		return c.Render(http.StatusOK, "goal", data)
	}
}

func handleEntry(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		goalId := c.FormValue("goalId")
		id, err := uuid.Parse(goalId)

		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}
		date := c.FormValue("date")
		parsed, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}
		parsed = utils.TruncateDate(parsed)
		complete, err := models.ToggleEntry(db, parsed, id)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		color, _ := models.FetchGoalColor(db, id)
		data := Entry{Date: parsed.Format(time.RFC3339), Complete: complete, Color: color}
		return c.Render(http.StatusOK, "entry", data)
	}
}

//go:generate npm run build
func main() {
	db := database.Init("db.db")
	database.Migrate(db)
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Logger())

	e.Static("static/css", "static/css")
	e.Static("static/js", "static/js")
	NewTemplateRenderer(e, "static/*.html")

	e.GET("/", func(c echo.Context) error {
		USER, _ := uuid.Parse("5c1c0569-dcd1-4c0d-87f0-0d0c1debdd5b")
		goals, _ := models.FetchGoals(db, USER)
		data := IndexTemplateData{Goals: mapGoals(goals)}
		return c.Render(http.StatusOK, "index", data)
	})
	e.POST("/entry", handleEntry(db))
	e.POST("/routine", handleRoutine(db))

	e.GET("/users/:userId", handlers.GetUser(db))
	e.POST("/users", handlers.PostUser(db))
	e.PATCH("/users/:userId", handlers.PatchUser(db))
	e.DELETE("/users/:userId", handlers.DeleteUser(db))

	e.GET("/goals", handlers.GetGoals(db))
	e.GET("/goals/:goalId", handlers.GetGoal(db))
	e.POST("/goals", handlers.PostGoal(db))
	e.PATCH("/goals/:goalId", handlers.PatchGoal(db))
	e.DELETE("/goals/:goalId", handlers.DeleteGoal(db))

	e.GET("/goals/:goalId/entries", handlers.GetEntries(db))
	e.POST("/goals/:goalId/entries", handlers.PutEntry(db))
	e.DELETE("/goals/:goalId/entries", handlers.DeleteEntry(db))
	e.Logger.Fatal(e.Start(":8000"))
}
