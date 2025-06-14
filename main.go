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
  home := os.Getenv("HOME_DIR")

  e.GET("/withdraw/export/:id", withdraw.ExportSetHandler)
  e.POST("/withdraw/process", withdraw.ProcessHandler)
  e.Static("/reports", "views/reports")
  e.File("/withdraw/set.html", home + "/views/withdraw/set.html")
  e.GET("/report_list", withdraw.ListReportsHandler)
  e.File("/version", home + "/version.txt")
  e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
