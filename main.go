package main

import (
	"io"
	"net/http"
	database "tb/goals/db"
	"tb/goals/handlers"
	"text/template"

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
	NewTemplateRenderer(e, "static/*.html")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/users/:userId", handlers.GetUser(db))
	e.POST("/users", handlers.PostUser(db))
	e.PATCH("/users/:userId", handlers.PatchUser(db))
	e.DELETE("/users/:userId", handlers.DeleteUser(db))

	e.GET("/goals", handlers.GetGoals(db))
	e.POST("/goals", handlers.PostGoal(db))
	e.PATCH("/goals/:goalId", handlers.PatchGoal(db))
	e.DELETE("/goals/:goalId", handlers.DeleteGoal(db))

	e.GET("/goals/:goalId/entries", handlers.GetEntries(db))
	e.PUT("/goals/:goalId/entries", handlers.PutEntry(db))
	e.DELETE("/goals/:goalId/entries", handlers.DeleteEntry(db))
	e.Logger.Fatal(e.Start(":8000"))
}
