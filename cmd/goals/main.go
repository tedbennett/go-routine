package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"net/http"
	"text/template"
	"time"
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

// Models
type Goal struct {
	Id      uuid.UUID   `json:"id"`
	Entries []GoalEntry `json:"entries"`
	Title   string      `json:"title"`
}

type GoalEntry struct {
	Date     time.Time `json:"date"`
	Complete bool      `json:"complete"`
}

type Date struct {
	year, day int
	month     time.Month
}

func (d *Date) fromTimeDate(year int, month time.Month, day int) {
	d.year, d.month, d.day = year, month, day
}

func (g *Goal) updateEntry() {
	today := Date{}
	today.year, today.month, today.day = time.Now().Date()
	for i, entry := range g.Entries {
		entryDate := Date{}
		entryDate.fromTimeDate(entry.Date.Date())
		if today == entryDate {
			g.Entries[i].Complete = !entry.Complete
			return
		}
	}
	// No entry found with today's date
	entry := GoalEntry{Date: time.Now(), Complete: true}
	g.Entries = append(g.Entries, entry)
}

var goals []Goal

func createGoal(title string) uuid.UUID {
	id := uuid.New()
	goal := Goal{Id: id, Entries: []GoalEntry{}, Title: title}
	goals = append(goals, goal)
	return id
}

func updateTodaysGoalEntry(id uuid.UUID) {
	for i := range goals {
		if goals[i].Id == id {
			goals[i].updateEntry()
			break
		}
	}
}

//go:generate npm run build
func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Logger())
	e.Static("static/css", "static/css")
	NewTemplateRenderer(e, "static/*.html")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/goals", func(c echo.Context) error {
		return c.JSON(http.StatusOK, goals)
	})

	e.POST("/goals", func(c echo.Context) error {
		title := c.FormValue("title")
		id := createGoal(title)
		return c.JSON(http.StatusAccepted, id)
	})

	e.GET("/goals/:goalId", func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			// Invalid uuid, won't exist
			return c.NoContent(http.StatusNotFound)
		}
		for _, goal := range goals {
			if goal.Id == id {
				return c.JSON(http.StatusOK, goal)
			}
		}
		return c.NoContent(http.StatusNotFound)
	})

	e.POST("/goals/:goalId/entries", func(c echo.Context) error {
		goalId := c.Param("goalId")
		id, err := uuid.Parse(goalId)
		if err != nil {
			// Invalid uuid, won't exist
			return c.NoContent(http.StatusNotFound)
		}
		updateTodaysGoalEntry(id)
		return c.NoContent(http.StatusAccepted)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
