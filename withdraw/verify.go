package withdraw

import (
  "github.com/labstack/echo/v4"
  "log"
  "fmt"
  "net/http"
  "rds_alma_tools/file"
  faktory "github.com/contribsys/faktory/client"
  "rds_alma_tools/connect"
  "rds_alma_tools/oclc"
  "os"
  "io"
  "bytes"
  "encoding/json"
)

func ExportVerifyHandler(c echo.Context)error{
  f, _ := c.FormFile("file")
  src, err := f.Open()
  bytedata,_ := io.ReadAll(src)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Unable to open file") }
  defer src.Close()
  //generate a filename to use throughout
  var fname interface{} = file.Filename()
  var worker = c.FormValue("worker")

  var stringdata interface{} = string(bytedata)

  client, err := faktory.Open()
  if err != nil{ log.Println(err); return c.String(http.StatusInternalServerError, err.Error())}
  //arg0 jobname, arg1 args for the job
  job := faktory.NewJob("VerifyJob", fname, stringdata)
  job.Queue = fmt.Sprintf("process%s", worker)
  job.ReserveFor = 7200
  retries := 0
  job.Retry = &retries
  err = client.Push(job)
  if err != nil{ log.Println(err); return c.String(http.StatusInternalServerError, err.Error()) }
  base_url := os.Getenv("HOME_URL")
  return c.HTML(http.StatusOK, fmt.Sprintf("<p>Export will be written to <a href=\"%s/reports/%s\">%s</a></p>", base_url, fname, fname))
}

func VerifyList(filename string, data []byte){
  lines := bytes.Split(data, []byte("\n"))
  datalist := []string{}
  for _, line := range lines{
    if string(line) == "" { break }
    old_data := LineMap(string(line))
    url := BuildItemLink(old_data["mms_id"], old_data["holding_id"], old_data["pid"])
    new_data := VerifyItem(url, old_data["oclc"])
    datalist = append(datalist, FormatLine(old_data) + "\t" + new_data)
  }
  file.WriteReport(filename, datalist)
}

func FormatLine(lineMap map[string]string) string{
  return fmt.Sprintf("%s\t%s\t%s\t%s", lineMap["mms_id"], lineMap["barcode"], lineMap["pid"], lineMap["location"])
}

func VerifyItem(url string, oclc_num string)(string){
  params := []string{"view=brief", "apikey=" + os.Getenv("ALMA_KEY")}
  data, err := connect.Get(url, params)
  if err != nil { log.Println(err); return "" }
  var r Record
  err = json.Unmarshal(data, &r)
  if err != nil { log.Println(err); return ""}
  oclc_is_set, err := oclc.CheckHolding(oclc_num)
  if err != nil { oclc_is_set = "unable to confirm" }
  return fmt.Sprintf("%s\t%s\t%s\t%s\t%s",r.Item_data.Location.Value, r.Item_data.Base_status.Desc,r.Bib_data.Bib_suppress, AllianceLinked(r.Bib_data.Network), oclc_is_set) 
}
