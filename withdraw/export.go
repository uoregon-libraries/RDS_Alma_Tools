package withdraw

import(
  "github.com/labstack/echo/v4"
  "log"
  "io"
  "net/http"
  "rds_alma_tools/connect"
)

// holds all data for return to user
type Export struct{
  brief_items []BriefItem
}

// only fields of interest for withdraw purposes
type BriefItem struct{
  mms_id string
  title string
  holding_id string
  item_pid string
  barcode string
  base_status string
  library string
  location string
  public_note string
  fulfillment_note string
  internal_note_1 string
  internal_note_2 string
  internal_note_3 string
}

func ExportSetHandler(c echo.Context)(string, error) {
  // get set
  param_ids := []string{ c.Param("id") }
  json, err := connect.Get("sets", param_ids)
  if err != nil {}
  // get links for items
  links, err := ExtractLinks(json)
  if err != nil {}
  var export Export
  // iterate through items 
  for _, link := range links {
    export.processItem(link)
  }

  // return export file using json stream option?
  c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
  c.Response().WriteHeader(http.StatusOK)
  return json.NewEncoder(c.Response()).Encode(export)
}

// returns list of links pulled from members using gjson?
func extractLinks(json string)([]string, error){}

// called during iteration through links to items
// creates BriefItems and adds them to the export
func (e Export)processItem(url string){}

// pulls fields of interest using gjson?
func (b BriefItem)buildBriefItem(json string)error{}

