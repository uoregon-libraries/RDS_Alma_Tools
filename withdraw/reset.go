package withdraw

import(
  "github.com/labstack/echo/v4"
  "github.com/tidwall/sjson"
  "github.com/tidwall/gjson"
  "rds_alma_tools/connect"
  "rds_alma_tools/file"
  "os"
  "io"
  "log"
  "fmt"
  "bytes"
  "time"
  "net/http"
)

func ResetHandler(c echo.Context)(error){
  f, _ := c.FormFile("file")
  src, err := f.Open()
  data,_ := io.ReadAll(src)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()
  //generate a filename to write log-type information to for the user
  filename := file.Filename()
  ProcessReset(filename, data)
  return c.String(http.StatusOK, fmt.Sprintf("Errors will be written to \"%s\"", filename))
}

func ProcessReset(filename string, data []byte){
  pids := ResetItems(filename, data)
  bibs := UniqueBibs(data)
  ResetStatus(filename, pids, bibs)
}

func ResetItem(data string)([]byte, connect.Response){
  lineMap := LineMap(data)
  url := BuildItemLink(lineMap["mms_id"], lineMap["holding_id"], lineMap["pid"])
  params := []string{ ApiKey() }
  itemRec, err := connect.Get(url, params)
  if err != nil { 
    if itemRec != nil { return nil, connect.Response{ Id: url, Message: connect.ExtractAlmaError(string(itemRec))} } else { return nil, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error())} } }

  newLocVal := lineMap["location"]
  newLibVal := lineMap["library"]
  internalNote3 := lineMap["internal_note_3"]
  //using sjson insert new library, location, append note
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.location.value", newLocVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.internal_note_3", internalNote3)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.value", newLibVal)
  itemRec,_ = sjson.DeleteBytes(itemRec, "bib_data")
  return itemRec, connect.Response{Id:"", Message: connect.BuildMessage("")}
}

func ResetItems(filename string, data []byte)(map[string]Eligible){
  debug := os.Getenv("DEBUG")
  var r connect.Report
  pids := map[string]Eligible{}
  lines := bytes.Split(data, []byte("\n"))
  for _, line := range lines{
    if string(line) == "" { break }
    itemRec, response := ResetItem(string(line))
    if response.Id != ""{
      r.Responses = append(r.Responses, response)
      continue
    }
    params := []string{ ApiKey() }
    url := gjson.GetBytes(itemRec, "link").String()
    if debug == "true" { url = os.Getenv("TEST_URL") }
    body, err := connect.Put(url, params, string(itemRec))
    if err != nil {
      if body != nil {
      r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.ExtractAlmaError(err.Error()) }) } else { 
      r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error()) }) }
    }
    if body != nil { r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.BuildMessage("success") } )
      pid := ExtractPid(url)
      lm := LineMap(string(line))
      pids[pid] = Eligible{ Locations: []string{ lm["location"] } }
    }
  }
  r.WriteReport(filename)
  return pids
}

func ResetStatus(filename string, list map[string]Eligible, eligibleLists map[string]Eligible){
  setid := os.Getenv("RESET_ITEM_STATUS_SET")
  jobid := os.Getenv("UPDATE_ITEM_STATUS_JOB_ID")
  err := UpdateSet("RESET_ITEM_STATUS_SET", "ITEM", list)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return}
    var params = []Param{
      Param{ Name: Val{ Value: "MISSING_STATUS_selected" }, Value: "true"},
      Param{ Name: Val{ Value: "MISSING_STATUS_value" }, Value: "NULL" },
      Param{ Name: Val{ Value: "MISSING_STATUS_condition" }, Value: "NULL" },
      Param{ Name: Val{ Value: "set_id" }, Value: setid },
    }
  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return}
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)
  CheckJob(instance, ResetUnsuppress, filename, eligibleLists)
}

func ResetNetworkLink(filename string, list map[string]Eligible){
  setid := os.Getenv("RESET_LINK_SET")
  jobid := os.Getenv("RESET_LINK_JOB_ID")

  err := UpdateSet("RESET_LINK_SET", "BIB_MMS", list)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  params := []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: setid },
    Param{ Name: Val{ Value: "contribute_nz" }, Value: "true" },
    Param{ Name: Val{ Value: "non_serial_match_profile" }, Value: "com.exlibris.repository.mms.match.uniqueOCLC" },
    Param{ Name: Val{ Value: "non_serial_match_prefix" }, Value: "" },
    Param{ Name: Val{ Value: "serial_match_profile" }, Value: "com.exlibris.repository.mms.match.uniqueOCLC" },
    Param{ Name: Val{ Value: "serial_match_prefix" }, Value: "" },
    Param{ Name: Val{ Value: "ignoreResourceType" }, Value: "false" },
  }

  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  CheckJob(instance, nil, filename, nil)
}

func ResetUnsuppress(filename string, list map[string]Eligible){
  setid := os.Getenv("UNSUPPRESS_SET")
  jobid := os.Getenv("SUPPRESS_JOB_ID")
  err := UpdateSet("UNSUPPRESS_SET", "BIB_MMS", list)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  params := []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: setid },
    Param{ Name: Val{ Value: "task_MmsTaggingParams_boolean" }, Value: "false" },
  }
  instance,err := SubmitJob(jobid, params)
  if err != nil { file.WriteReport(filename, []string{ err.Error() }); return }
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  time.Sleep(span)

  CheckJob(instance, ResetNetworkLink, filename, list)
  }
