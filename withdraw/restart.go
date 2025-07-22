package withdraw

import (
  "github.com/labstack/echo/v4"
  "log"
  "rds_alma_tools/connect"
  faktory "github.com/contribsys/faktory/client"
  "os"
  "io"
  "net/http"
  "rds_alma_tools/file"
  "encoding/json"
  "fmt"
  "bytes"
)


func RestartHandler(c echo.Context)(error){
  //get uploaded file
  f, _ := c.FormFile("file")
  src, err := f.Open()
  bytedata,_ := io.ReadAll(src)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()
  //generate a filename to use throughout
  var fname interface{} = file.Filename()
  var worker = c.FormValue("worker")
  var stage interface{} = c.FormValue("stage")

  var stringdata interface{} = string(bytedata)

  client, err := faktory.Open()
  if err != nil{ log.Println(err); return c.String(http.StatusInternalServerError, err.Error())}
  //arg0 jobname, arg1 args for the job
  job := faktory.NewJob("RestartJob", fname, stage, stringdata)
  job.Queue = fmt.Sprintf("process%s", worker)
  job.ReserveFor = 7200
  retries := 0
  job.Retry = &retries
  err = client.Push(job)
  if err != nil{ log.Println(err); return c.String(http.StatusInternalServerError, err.Error()) }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Relevant updates will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}

//this is analagous to Process
func Restart(fname, stage string, data []byte){
  eligibleLists, errs := EligibleToUnlinkSuppressUnsetList(data)
  if len(errs) != 0 { file.WriteReport(fname, errs)}
  if stage == "status" {
    pids := CollectUpdatedItems(fname, data)
    ProcessStatusUpdate(fname, pids, eligibleLists)
  } else if stage == "unlink" { 
    ProcessUnlink(fname, eligibleLists) 
  } else if stage == "suppress" {
    ProcessSuppress(fname, eligibleLists)
  } else {
    ProcessUnset(fname, eligibleLists)
  }
}

func CollectUpdatedItems(fname string, data []byte)(map[string]Eligible){
  lines := bytes.Split(data, []byte("\n"))
  pids := map[string]Eligible{}
  var r connect.Report
  for _, line := range lines{
    if string(line) == "" { break }
    old_data := LineMap(string(line))
    url := BuildItemLink(old_data["mms_id"], old_data["holding_id"], old_data["pid"])
    missing, err := MissingStatus(url)
    if err != nil { r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error()) })}
    if !missing {
      pids[old_data["pid"]] = Eligible{ Locations: []string{ old_data["location"] } }
    }
  }
  r.WriteReport(fname)
  return pids
}

func MissingStatus(url string)(bool, error){
  params := []string{"view=brief", "apikey=" + os.Getenv("ALMA_KEY")}
  data, err := connect.Get(url, params)
  if err != nil { log.Println(err); return false, err }
  var r Record
  err = json.Unmarshal(data, &r)
  if err != nil { log.Println(err); return false, err}
  if r.Item_data.Library.Value == "Withdrawn" && r.Item_data.Base_status.Desc == "Item in place" {
    return false, nil
  } else {
    return true, nil
  }
}
