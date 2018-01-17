package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olebedev/config"
)

func main() {
	// Parse config yaml string from ./conf.go
	conf, err := config.ParseYaml(confString)
	PanicIf(err)

	// Set config variables delivered from main.go:11
	// Variables defined as ./conf.go:3
	conf.Set("debug", debug)
	conf.Set("commitHash", commitHash)

	// Parse environ variables from system environment
	conf.Env()

	// Echo instance
	engine := echo.New()

	// Set up echo debug level
	engine.Debug = conf.UBool("debug") // get debug config value returns boolean

	// Middleware
	// config logger
	engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${method} | ${status} | ${uri} -> ${latency_human}` + "\n",
	}))
	engine.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	engine.Renderer = t

	// Route => handler
	engine.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	engine.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "hello", "world")
	})

	// Start server
	engine.Logger.Fatal(engine.Start(":" + conf.UString("port")))
}
