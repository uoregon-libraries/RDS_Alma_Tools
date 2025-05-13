package withdraw

import(
  "github.com/labstack/echo/v4"
  "rds_alma_tools/oclc"
  "rds_alma_tools/file"
  "log"
  "net/http"
  "os"
  "net/url"
  "time"
  "io"
  "fmt"
)

func ProcessHandler(c echo.Context)(error){
  f, _ := c.FormFile("file")
  src, err := f.Open()
  data,_ := io.ReadAll(src)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()
  //generate a filename to write log-type information to for the user
  filename := file.Filename()

  loc_type := c.FormValue("loc_type")
  if loc_type == "" { return c.String(http.StatusBadRequest, "Location type is required") }
  Process(filename, loc_type, data)
  return c.String(http.StatusOK, fmt.Sprintf("Relevant updates will be written to \"%s\"", filename))
}

func Process(filename, loc_type string, data []byte){
  pids := UpdateItems(filename, loc_type, data)
  eligibleLists, errs := EligibleToUnlinkSuppressUnsetList(data)

  if len(errs) != 0 { file.WriteReport(filename, errs)}
  ProcessStatusUpdate(filename, pids, eligibleLists)
}

func ProcessStatusUpdate(filename string, list map[string]Eligible, eligibleLists map[string]Eligible ){
  setid := os.Getenv("UPDATE_ITEM_STATUS_SET")
  jobid := os.Getenv("UPDATE_ITEM_STATUS_JOB_ID")
  err := UpdateSet("UPDATE_ITEM_STATUS_SET", "ITEM", list)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return}
    var params = []Param{
      Param{ Name: Val{ Value: "MISSING_STATUS_selected" }, Value: "true"},
      Param{ Name: Val{ Value: "MISSING_STATUS_value" }, Value: "MISSING" },
      Param{ Name: Val{ Value: "MISSING_STATUS_condition" }, Value: "NULL" },
      Param{ Name: Val{ Value: "set_id" }, Value: setid },
    }
  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return}
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)
  CheckJob(instance, ProcessUnlink, filename, eligibleLists)
}

func ProcessUnlink(filename string, eligibleLists map[string]Eligible){
  setid := os.Getenv("UNLINK_SET")
  jobid := os.Getenv("UNLINK_JOB_ID")
  unlinkList := Winnow(eligibleLists, OkToUnlink)
  if len(unlinkList) == 0 { 
    file.WriteReport(filename, []string{ "Nothing to unlink" })
    ProcessSuppress(filename, eligibleLists)
    return
  }
  err := UpdateSet("UNLINK_SET", "BIB_MMS", unlinkList)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  params := []Param{ Param{ Name: Val{ Value: "set_id" }, Value: setid } }

  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  CheckJob(instance, ProcessSuppress, filename, eligibleLists)
}

func ProcessSuppress(filename string, eligibleLists map[string]Eligible){
  setid := os.Getenv("SUPPRESS_SET")
  jobid := os.Getenv("SUPPRESS_JOB_ID")
  suppressList := Winnow(eligibleLists, OkToSuppress)
  if len(suppressList) == 0 {
    file.WriteReport(filename, []string{ "Nothing to suppress" })
    ProcessUnset(filename, eligibleLists)
    return
  }
  err := UpdateSet("SUPPRESS_SET", "BIB_MMS", suppressList)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  params := []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: setid },
    Param{ Name: Val{ Value: "task_MmsTaggingParams_boolean" }, Value: "true" },
  }

  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }

  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  CheckJob(instance, ProcessUnset, filename, eligibleLists)
}

func ProcessUnset(filename string, eligibleLists map[string]Eligible){
  unsetlist := Winnow(eligibleLists, OkToUnset)
  newlist := map[string]string{}
  for k,v := range unsetlist{
    newlist[k] = v.Oclc
  }
  if len(newlist) == 0 {
    file.WriteReport(filename, []string{ "Nothing to unset" })
    return
  }
  oclc.UnsetHoldings(filename, newlist)
  Final_Report(filename, eligibleLists)
}

func Final_Report(filename string, eligibleLists map[string]Eligible){
  results := []string{}
  for k, v := range eligibleLists{
    if v.SerialRequiresAction { results = append(results, k + ": " + "Serial Requires Action") }
  }
  file.WriteReport(filename, results)
}

func BuildItemLink(mmsId string, holdingId string, pid string)string{
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("bibs", mmsId, "holdings", holdingId, "items", pid)
  return _url.String()
}

func BuildBibLink(mmsId string)string{
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("bibs", mmsId)
  return _url.String()
}

func BuildHoldingLink(mmsId string, holdingId string)string{
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("bibs", mmsId, "holdings", holdingId)
  return _url.String()

}

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}

func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}

type Selector func(e Eligible) bool

func OkToUnlink(e Eligible) bool{
  return e.Unlink
}

func OkToSuppress(e Eligible) bool{
  return e.Suppress
}

func OkToUnset(e Eligible) bool{
  return e.Unset
}

func Winnow(list map[string]Eligible, sfunc Selector) map[string]Eligible{
  newlist := map[string]Eligible{}
  for k, v:= range list{
    if sfunc(v) == true {
      newlist[k] = v
    }
  }
  return newlist
}
