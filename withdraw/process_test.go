package withdraw

import(
  "testing"
  "github.com/tidwall/gjson"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "time"
  //"strings"
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
    fmt.Fprintf(w, string(data))
   }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/")

  fy := FiscalYear(TimeNow())
  line := "9984898401852\tXBox 360\t22274069860001852\t23193212440001852\t35025040997286\tItem not in place\tScience\tsgames\tfake public note\ttoggled missing status from technical migration. was breaking bookings - SDG\tSTATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675\tNOTE(ITEM): serial number: 118381693005\tStatus: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG\tfake_retention_note\n"

  itemRec, _ := UpdateItem("withdraw", line)
  library := gjson.GetBytes(itemRec, "item_data.library.value")
  if library.String() != "swithdrwn" { t.Errorf("library is wrong") }
  inote3 := gjson.GetBytes(itemRec, "item_data.internal_note_3")
  if inote3.String() != "Status: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG|WD FY" + fy {
    t.Errorf("internal note 3 is wrong")
  }
  status := gjson.GetBytes(itemRec, "item_data.base_status.desc")
  if status.String() != "missing" { t.Errorf("status is wrong") }
}

func TestCheck_library(t *testing.T){
}

func TestHolding_items(t *testing.T){
}

func TestBib_holdings(t *testing.T){
}
