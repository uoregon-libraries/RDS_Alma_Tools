package main

import (
  "rds_alma_tools/withdraw"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "os"
)

func main() {
  e := echo.New()
  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  e.GET("/withdraw/export/:id", withdraw.ExportSetHandler)
  e.POST("/withdraw/process", withdraw.ProcessHandler)
  e.Static("/withdraw", "views/withdraw") //urlpath,directorypath, withdraw/set.html
  e.Static("/reports", "views/reports")
  e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
