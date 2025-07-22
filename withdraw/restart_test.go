package withdraw

import (
  "testing"
  "net/http/httptest"
  "fmt"
  "net/http"
  "os"
)

func TestMissingStatus(t *testing.T){
  data1 := `{"item_data": { "library": { "value": "Withdrawn", "desc": "Withdrawn Library" }, "location": { "value": "kwithdrwn", "desc": "Knight withdrawn" }, "base_status": { "desc": "Item in place" } } }`
  data2 := `{"item_data": { "library": { "value": "Withdrawn", "desc": "Withdrawn Library" }, "location": { "value": "kwithdrawn", "desc": "Knight withdrawn" }, "base_status": { "description": "Item not in place" } } }`

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/bibs/123/holdings/456/items/7890" {
      fmt.Fprintf(w, data1)
    } else if r.URL.Path == "/bibs/456/holdings/789/items/1230" {
      fmt.Fprintf(w, data2)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")

  missing,_ := MissingStatus(ts.URL + "/bibs/123/holdings/456/items/7890")
  if missing != false { t.Errorf("incorrect missing status") }
  missing,_ = MissingStatus(ts.URL + "/bibs/456/holdings/789/items/1230")
  if missing != true { t.Errorf("incorrect missing status") }

}
