package withdraw

import(
  "github.com/labstack/echo/v4"
  //"rds_alma_tools/oclc"
  "log"
  "net/http"
  "os"
  "net/url"
  "time"
  "io"
  "fmt"
)

func ProcessHandler(c echo.Context)(error){
  //get uploaded file
  file, _ := c.FormFile("file")
  src, err := file.Open()
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()
  //generate a filename to use throughout
  filename := Filename()

  loc_type := c.FormValue("loc_type")
  if loc_type == "" { return c.String(http.StatusBadRequest, "Location type is required") }
  Process(filename, loc_type, src)
  return c.String(http.StatusOK, fmt.Sprintf("Relevant updates will be written to \"%s\"", filename))
}

// runs initial steps and then launches the series of steps that require waiting
func Process(filename, loc_type string, src io.Reader){
  pids := UpdateItems(filename, loc_type, src) // THIS ONLY RETURNS SUCCESSFUL ITEMS?
  var eligibleLists = map[string][]bool{}
  eligibleLists, err := EligibleToUnlinkAndSuppressList(src) //THIS WORKS FROM ORIG LIST, NOT ITEMS
  if len(eligibleLists) == 0 { log.Println("eligibleLists starts empty"); return }
  if err != nil { WriteReport(filename, err.Error())}// CONTINUE ANYWAY
  ProcessStatusUpdate(filename, pids, eligibleLists)
}

// item updates that require a job
func ProcessStatusUpdate(filename string, list []string, eligibleLists map[string][]bool ){
  setid := os.Getenv("UPDATE_ITEM_STATUS_SET")
  jobid := os.Getenv("UPDATE_ITEM_STATUS_JOB_ID")

  err := UpdateSet(filename, "UPDATE_ITEM_STATUS_SET", "ITEM", list)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return}

var params = []Param{
    Param{ Name: Val{ Value: "MISSING_STATUS_selected" }, Value: "true"},
    Param{ Name: Val{ Value: "MISSING_STATUS_value" }, Value: "MISSING" },
    Param{ Name: Val{ Value: "MISSING_STATUS_condition" }, Value: "NULL" },
    Param{ Name: Val{ Value: "set_id" }, Value: setid }, 
  }

  instance,err := SubmitJob(filename, jobid, params)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return}

  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  if len(eligibleLists) == 0 {
    log.Println("eligibleLists empty")
    CheckJob(instance, nil, filename, nil)
  } else { CheckJob(instance, ProcessUnlink, filename, eligibleLists) }
}

func ProcessUnlink(filename string, eligibleLists map[string][]bool){
  setid := os.Getenv("UNLINK_SET")
  jobid := os.Getenv("UNLINK_JOB_ID")
  unlinkList := extractEligibles(eligibleLists, 0)

  err := UpdateSet(filename, "UNLINK_SET", "BIB_MMS", unlinkList)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return }
  params := []Param{ Param{ Name: Val{ Value: "set_id" }, Value: setid } }

  instance,err := SubmitJob(filename, jobid, params)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return }
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  CheckJob(instance, ProcessSuppress, filename, eligibleLists)
}

// currently, this will not be called unless the unlink job returns COMPLETED_SUCCESS
// TBD how to continue if some items can/should be suppressed
func ProcessSuppress(filename string, eligibleLists map[string][]bool){
  setid := os.Getenv("SUPPRESS_SET")
  jobid := os.Getenv("SUPPRESS_JOB_ID")
  suppressList := extractEligibles(eligibleLists, 1)

  err := UpdateSet(filename, "SUPPRESS_SET", "BIB_MMS", suppressList)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return }
  params := []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: setid },
    Param{ Name: Val{ Value: "task_MmsTaggingParams_boolean" }, Value: "true" },
  }

  instance,err := SubmitJob(filename, jobid, params)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return }

  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  // eventually call oclc.Remove
  CheckJob(instance, nil, filename, eligibleLists)
}

func extractEligibles(lists map[string][]bool, ind int)[]string{
  newlist := []string{}
  for k,v := range lists{
    if v[ind] { newlist = append(newlist, k) }
  }
  return newlist
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

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}

func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}
