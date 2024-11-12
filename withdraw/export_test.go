package withdraw

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "strings"
)

func TestExtractLinks(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, _ := os.ReadFile(homedir + "/fixtures/set_members.json")
  links := ExtractLinks(data)
  if !Contains(links, "https://api-na.hosted.exlibrisgroup.com/almaws/v1/bibs/9994938301852/holdings/22348542010001852/items/23189116110001852") {  t.Fatalf("link is not present") }
  if !Contains(links,  "https://api-na.hosted.exlibrisgroup.com/almaws/v1/bibs/9994936001852/holdings/22189197510001852/items/23189197500001852") {  t.Fatalf("link is not present") }
}

func TestProcessItem(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, _ := os.ReadFile(homedir + "/fixtures/item_record.json")
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, string(data))
   }))
  defer ts.Close()

  str := ProcessItem(ts.URL)
  arr := strings.Split(str, "\t")
  if arr[0] != "9984898401852" { t.Fatalf("mms_id incorrect") }
  if arr[1] != "XBox 360" { t.Fatalf("title incorrect") }
  if arr[2] != "22274069860001852" { t.Fatalf("holding id incorrect") }
  if arr[3] != "23193212440001852" { t.Fatalf("pid incorrect") }
  if arr[4] != "35025040997286" { t.Fatalf("barcode incorrect") }
  if arr[5] != "Item not in place" { t.Fatalf("base_status incorrect") }
  if arr[6] != "Science" { t.Fatalf("library incorrect") }
  if arr[7] != "sgames" { t.Fatalf("location incorrect") }
  if arr[8] != "fake public note" { t.Fatalf("public_note incorrect") }
  if arr[9] != "toggled missing status from technical migration. was breaking bookings - SDG" { t.Fatalf("fulfillment note incorrect") }
  if arr[10] != "STATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675" { t.Fatalf("internal note 1 incorrect") }
  if arr[11] != "NOTE(ITEM): serial number: 118381693005" { t.Fatalf("internal note 2 incorrect") }
  if arr[12] != "Status: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG" { t.Fatalf("internal note 3 incorrect") }
}

func Contains(hay []string, needle string) bool {
  for _, value := range hay {
    if value == needle {
      return true
    }
  }
  return false
}
