package withdraw

import(
  "github.com/labstack/echo/v4"
  "os"
  "log"
  "net/http"
)

func ListReportsHandler(c echo.Context)error{
  entries, err := os.ReadDir(os.Getenv("REPORT_DIR"))
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, err.Error()) }
  list := ""
  for _, e := range entries {
    list += e.Name() + "\n"
  }
  return c.String(http.StatusOK, list)
}
