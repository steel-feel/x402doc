package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v5"
)

func main() {
	PORT := flag.String("port", "9000", "port that will be exposed")
	flag.Parse()

	e := echo.New()
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]any{
			"message": "Hello, World!",
		})
	})


	
	if err := e.Start(fmt.Sprintf(":%s", *PORT)); err != nil {
    e.Logger.Error("failed to start server", "error", err)
  }
}
