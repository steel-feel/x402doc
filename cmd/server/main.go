package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v5"
	health "github.com/steel-feel/prac/internal/app/health"
)

func main() {
	PORT := flag.String("port", "9000", "port that will be exposed")
	flag.Parse()

	e := echo.New()

	healthHandler := health.NewHandler()
	baseRoute := e.Group("/api/v1")

	healthHandler.RegisterRoutes(baseRoute)

	if err := e.Start(fmt.Sprintf(":%s", *PORT)); err != nil {
    e.Logger.Error("failed to start server", "error", err)
  }
}
