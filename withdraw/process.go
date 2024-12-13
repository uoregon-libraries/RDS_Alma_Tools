package withdraw

import(
  "github.com/labstack/echo/v4"
  "rds_alma_tools/connect"
  "log"
  "net/http"
  "os"
  "net/url"
  "strconv"
  "time"
)

func ProcessHandler(c echo.Context)(error){
  //get uploaded file
  file, _ := c.FormFile("file")
  src, err := file.Open()
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()

  var report connect.Report
  loc_type := c.Param("loc_type")
  report = UpdateItems(report, loc_type, src)
  eligibleList, err := EligibleToUnlinkAndSuppressList(src)
  if err != nil { return c.String(http.StatusBadRequest, report.ResponsesToString()) }

  if len(eligibleList) == 0 { return c.String(http.StatusOK, report.ResponsesToString()) }

  //next steps...

  return c.String(http.StatusOK, report.ResponsesToString())
}

func BuildItemLink(mmsId string, holdingId string, pid string)string{
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("bibs", mmsId, "holdings", holdingId, "items", pid)
  return _url.String()
}

func TimeNow()time.Time{
  loc, _ := time.LoadLocation("America/Los_Angeles")
  t := time.Now().In(loc)
  return t
}

func FiscalYear(t time.Time)string{
  m := t.Format("01")
  y := t.Format("2006")
  mInt, _ := strconv.Atoi(m)
  if mInt > 6 { return y } else {
    yInt, _ := strconv.Atoi(y)
    return strconv.Itoa(yInt-1)
  }
}

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}

func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}
