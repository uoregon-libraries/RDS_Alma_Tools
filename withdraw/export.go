package withdraw

import(
  "github.com/labstack/echo/v4"
  "log"
  "rds_alma_tools/connect"
  "github.com/tidwall/gjson"
  "os"
  "net/url"
  "encoding/json"
  "strings"
)

func ExportSetHandler(c echo.Context)(error) {
  base_url := os.Getenv("ALMA_URL")
  full_url, _ := url.JoinPath(base_url, "conf", "sets", c.Param("id"), "members")
  params := []string{"limit=100", "apikey=" + os.Getenv("ALMA_KEY")}

  //get set
  set_data, err := connect.Get(full_url, params)
  if err != nil {log.Println(err); return echo.NewHTTPError(400, "Error retrieving set")}
  // get links for items
  links := ExtractLinks([]byte(set_data))

  //prepare file
  f, err := os.CreateTemp("","withdraw_set-")
  if err != nil {log.Println(err); return echo.NewHTTPError(400, "Error creating temporary file")}
  defer f.Close()
  defer os.Remove(f.Name())
  _, err = f.WriteString(BriefItemHead() + "\n")
  if err != nil { log.Println(err); return echo.NewHTTPError(400, "Error writing header")}

  // iterate through items 
  str := ""
  for _, link := range links {
    str = ProcessItem(link)
    _, err := f.WriteString(str + "\n")
    if err != nil {log.Println(err); return echo.NewHTTPError(400, "Error writing data")}
  }
  filename := "export_" + c.Param("id") + ".tsv"
  return c.Attachment(f.Name(), filename)
}

// returns list of links pulled from members using gjson
func ExtractLinks(json []byte)([]string){
  arr := []string{}
  members := gjson.GetBytes(json, "member.#.link")
  for _, link := range members.Array(){
    arr = append(arr, link.String())
  }
  return arr
}

func BriefItemHead()string{
  properties := []string{
    "mms_id",
    "title",
    "holding_id",
    "pid",
    "barcode",
    "base_status",
    "library",
    "location",
    "public_note",
    "fulfillment_note",
    "internal_note_1",
    "internal_note_2",
    "internal_note_3",
    "retention_note"}
  return strings.Join(properties[:], "\t")
}

// called during iteration through links to items
// creates BriefItems and adds them to the export
func ProcessItem(url string)string{
  params := []string{"view=brief", "apikey=" + os.Getenv("ALMA_KEY")}
  data, err := connect.Get(url, params)
  if err != nil { log.Println(err); return "" }
  var r Record
  err = json.Unmarshal([]byte(data), &r)
  if err != nil { log.Println(err); return ""}
  str := r.Stringify()
  return str
}

