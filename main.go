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
  //e.GET("/withdraw/process", withdraw.ProcessHandler)
  //e.Static("withdraw", "views") //urlpath,directorypath, withdraw/set.html
  e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
