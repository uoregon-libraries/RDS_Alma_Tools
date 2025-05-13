package withdraw

import (

  "rds_alma_tools/connect"
  "rds_alma_tools/file"
  "github.com/tidwall/sjson"
  "os"
  "github.com/tidwall/gjson"
  "strings"
  "io"
  "time"
  "strconv"
  "bytes"
)

//returns list of item pids
func UpdateItems(filename string, loc_type string, data []byte)(map[string]Eligible){
  debug := os.Getenv("DEBUG")
  var r connect.Report
  pids := map[string]Eligible{}
  lines := bytes.Split(data, []byte("\n"))
  for _, line := range lines{
    if string(line) == "" { break }
    itemRec, response := UpdateItem(loc_type, string(line))
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

func UpdateItem(newLocType string, data string)([]byte, connect.Response){
  lineMap := LineMap(data)
  url := BuildItemLink(lineMap["mms_id"], lineMap["holding_id"], lineMap["pid"])
  params := []string{ ApiKey() }
  itemRec, err := connect.Get(url, params)
  if err != nil { 
    if itemRec != nil { return nil, connect.Response{ Id: url, Message: connect.ExtractAlmaError(string(itemRec))} } else { return nil, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error())} } }
  libMap := WithdrawDeselectMap()
  newLocVal := libMap[WDKey{lineMap["library"], newLocType, "value"}]
  if newLocVal == "" { return nil, connect.Response{ Id: url, Message: connect.BuildMessage("Unable to determine new location") } }
  newLibVal := "Withdrawn"
  internalNote3 := lineMap["internal_note_3"] + "|WD FY" + FiscalYear(file.TimeNow())
  //using sjson insert new library, location, append note
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.location.value", newLocVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.internal_note_3", internalNote3)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.value", newLibVal)
  itemRec,_ = sjson.DeleteBytes(itemRec, "bib_data")
  return itemRec, connect.Response{Id:"", Message: connect.BuildMessage("")}
}

type WDKey struct{
  Libname, Libtype, Proptype string
}

func WithdrawDeselectMap()map[WDKey]string{
  libmap := map[WDKey]string{}
  homedir := os.Getenv("HOME_DIR")
  src,_ := os.Open(homedir + "/withdraw/library_map.txt")
  data,_ := io.ReadAll(src)
  lines := bytes.Split(data, []byte("\n"))
  for _,line := range lines{
    if string(line) == "" { break }
    arr := bytes.Split(line, []byte("\t"))
    libmap[WDKey{string(arr[0]), string(arr[1]), string(arr[2])}] = string(arr[3])
  }
  return libmap
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

func ExtractPid(url string) string{
  arr := strings.Split(url, "/")
  length := len(arr)
  return arr[length-1]
}
