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
  data, err := os.ReadFile(homedir + "/fixtures/set_members.json")
  if err != nil { t.Fatalf("did not read file") }
  links := ExtractLinks(data)
  if !Contains(links, "https://api-na.hosted.exlibrisgroup.com/almaws/v1/bibs/9994938301852/holdings/22348542010001852/items/23189116110001852") {  t.Errorf("first link is not present") }
  if !Contains(links,  "https://api-na.hosted.exlibrisgroup.com/almaws/v1/bibs/9994936001852/holdings/22189197510001852/items/23189197500001852") {  t.Errorf("second link is not present") }
}

func TestProcessItem(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, _ := os.ReadFile(homedir + "/fixtures/item_record.json")
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, string(data))
   }))
  defer ts.Close()

  str := ProcessItem(ts.URL)
  fmt.Println(str)
  arr := strings.Split(str, "\t")
  fmt.Println("length")
  fmt.Println(arr[14])
  if arr[0] != "9984898401852" { t.Errorf("mms_id incorrect") }
  if arr[1] != "XBox 360" { t.Errorf("title incorrect") }
  if arr[2] != "12345678" { t.Errorf("oclc incorrect") }
  if arr[3] != "22274069860001852" { t.Errorf("holding id incorrect") }
  if arr[4] != "23193212440001852" { t.Errorf("pid incorrect") }
  if arr[5] != "35025040997286" { t.Errorf("barcode incorrect") }
  if arr[6] != "Item not in place" { t.Errorf("base_status incorrect") }
  if arr[7] != "Science" { t.Errorf("library incorrect") }
  if arr[8] != "sgames" { t.Errorf("location incorrect") }
  if arr[9] != "fake public note" { t.Errorf("public_note incorrect") }
  if arr[10] != "toggled missing status from technical migration. was breaking bookings - SDG" { t.Errorf("fulfillment note incorrect") }
  if arr[11] != "STATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675" { t.Errorf("internal note 1 incorrect") }
  if arr[12] != "NOTE(ITEM): serial number: 118381693005" { t.Errorf("internal note 2 incorrect") }
  if arr[13] != "Status: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG" { t.Errorf("internal note 3 incorrect") }
  if arr[14] != "fake retention note" { t.Errorf("retention note incorrect") }
}

func Contains(hay []string, needle string) bool {
  for _, value := range hay {
    if value == needle {
      return true
    }
  }
  return false
}
