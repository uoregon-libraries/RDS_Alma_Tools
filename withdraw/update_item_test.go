package withdraw

import(
  "testing"
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "time"
  "rds_alma_tools/connect"
  "strings"
  "io"
)

func TestLoadMap(t *testing.T){
  lmap := LoadMap()
  if lmap[Key{"Knight","withdraw","value"}] != "kwithdrwn" {
    t.Errorf("new library value is wrong")
  }
}

func TestFiscalYear(t *testing.T){
  fakeTime, _ := time.Parse("2006-01 MST", "2024-11 PST")
  y := FiscalYear(fakeTime)
  if y != "2024" { t.Errorf("fiscal year is wrong") }
}

func TestUpdateItem(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, err := os.ReadFile(homedir + "/fixtures/item_record.json")
  if err != nil { t.Fatalf("did not read file") }
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/almaws/v1/bibs/9984898401852/holdings/22274069860001852/items/23193212440001852" { t.Errorf("wrong request url") }
    fmt.Fprintf(w, string(data))
   }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")

  fy := FiscalYear(TimeNow())
  line := "9984898401852\tXBox 360\t12345678\t22274069860001852\t23193212440001852\t35025040997286\tItem not in place\tScience\tsgames\tfake public note\ttoggled missing status from technical migration. was breaking bookings - SDG\tSTATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675\tNOTE(ITEM): serial number: 118381693005\tStatus: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG\tfake_retention_note\n"

  itemRec, _ := UpdateItem("withdraw", line)
  location := gjson.GetBytes(itemRec, "item_data.location.value")
  if location.String() != "swithdrwn" { t.Errorf("location is wrong") }
  inote3 := gjson.GetBytes(itemRec, "item_data.internal_note_3")
  if inote3.String() != "Status: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG|WD FY" + fy {
    t.Errorf("internal note 3 is wrong")
  }
  status := gjson.GetBytes(itemRec, "item_data.library.value")
  if status.String() != "Withdrawn" { t.Errorf("library is wrong") }
}

func TestSJSON(t *testing.T){
  data := `{"item_data":{"library": {"value": "blah"}}}`
  data2 := `{location: {"value": "blah"} }`
  item, err := sjson.Set(data, "item_data.library.value", "newval")
  if err != nil { t.Errorf(err.Error()) }
  if item != `{"item_data":{"library": {"value": "newval"}}}` { t.Errorf("wrong val") }
  item, err = sjson.Set(data2, "item_data.library.value", "newval")
  if err != nil { t.Errorf(err.Error()) } else {
    fmt.Println("sjson just adds the missing elts if needed, does not throw an error: " + item)
  }
}

func TestUpdateItems(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, err := os.ReadFile(homedir + "/fixtures/item_record.json")
  if err != nil { t.Fatalf("did not read file") }
  path := "/almaws/v1/bibs/9984898401852/holdings/22274069860001852/items/23193212440001852"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      if r.URL.Path != path { t.Errorf("wrong request url") }
      fmt.Fprintf(w, string(data))
    } else {
      body, _ := io.ReadAll(r.Body)
      fmt.Fprintf(w, string(body))
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("ALMA_KEY", "123456789")
  os.Setenv("DEBUG", "true")
  os.Setenv("VERBOSE", "true")
  os.Setenv("TEST_URL", ts.URL + path)
  var report connect.Report
  line := "9984898401852\tXBox 360\t12345678\t22274069860001852\t23193212440001852\t35025040997286\tItem not in place\tScience\tsgames\tfake public note\ttoggled missing status from technical migration. was breaking bookings - SDG\tSTATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675\tNOTE(ITEM): serial number: 118381693005\tStatus: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG\tfake_retention_note\n"

  src := strings.NewReader(line)
  report = UpdateItems(report, "withdraw", src)
  id := gjson.Get(report.ResponsesToString(), "id")
  if !strings.Contains(id.String(), path) { t.Errorf("response does not contain id") }
  message := gjson.Get(report.ResponsesToString(), "report.message")
  if message.String() != "success" { t.Errorf("response does not contain message") }

}
