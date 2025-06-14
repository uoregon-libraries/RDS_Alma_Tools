package withdraw

import(
  "testing"
  "net/http/httptest"
  "fmt"
  "net/http"
  "os"
  "slices"
  "log"
  "io"
)

func TestUniqueBibs(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  src, err := os.Open(homedir + "/fixtures/export.tsv")
  data,_ := io.ReadAll(src)
  if err != nil { t.Fatalf("did not read file") }
  bibs := UniqueBibs(data)
  if len(bibs) != 2 { t.Errorf("should be size 2") }
  _, ok := bibs["9984898401852"]
  if !ok { t.Errorf("bib not included") }
  _, ok = bibs["9984898401853"]
  if !ok { t.Errorf("bib not included") }
  for _, v := range bibs["9984898401852"].Locations{
    if v != "sgames" { t.Errorf("incorrect location slice" ) }
  }
}

func TestBibItems(t *testing.T){
  link1 := "https://api-na.hosted.exlibrisgroup.com/banana"
  link2 := "https://api-na.hosted.exlibrisgroup.com/cherry"
  data := `{"item":[{"link":"` + link1 + `"},
  {"link":"` + link2 + `"}]}`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/bibs/123456789/holdings/all/items" { t.Errorf("wrong request url") }
    log.Println(r.URL.Path)
    log.Println(data)
    fmt.Fprintf(w, data)
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")
  links, _ := BibItems("123456789")
  if !slices.Contains(links, link1) { t.Errorf("link is missing") }
  if !slices.Contains(links, link2) { t.Errorf("link is missing") }
}

func TestItemLibraryLocation(t *testing.T){
  data1 := `{"item_data": { "library": { "value": "Withdrawn", "desc": "Withdrawn Library" }, "location": { "value": "kwithdrwn", "desc": "Knight withdrawn" } } }`
  data2 := `{"item_data": { "library": { "value": "Science", "desc": "Price Science Commons" }, "location": { "value": "swithdrwn", "desc": "Science withdrawn" } } }`
  data3 := `{"item_data": { "library": { "value": "Department", "desc": "UO Departmental Library" }, "location": { "value": "zartmus", "desc": "Music" } } }`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/bibs/123/holdings/456/items/7890" {
      fmt.Fprintf(w, data1)
    } else if r.URL.Path == "/bibs/456/holdings/789/items/1230" {
      fmt.Fprintf(w, data2)
    } else {
      fmt.Fprintf(w, data3)
    }
   }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")

  result, _ := ItemLibraryLocation(ts.URL + "/bibs/123/holdings/456/items/7890")
  if result.LocCode != "kwithdrwn" {t.Errorf("wrong location")}

  result, _ = ItemLibraryLocation(ts.URL + "/bibs/456/holdings/789/items/1230")
  if result.LocCode != "swithdrwn" {t.Errorf("wrong location")}

  result, _ = ItemLibraryLocation(ts.URL + "/bibs/789/holdings/123/items/4560")
  if result.LocCode != "zartmus" {t.Errorf("wrong location")}
}

func TestEligibleToUnlinkSuppressUnset(t *testing.T){
  data1 := `{"item_data": { "library": { "value": "Withdrawn", "desc": "Withdrawn Library" }, "location": { "value": "kwithdrwn", "desc": "Knight withdrawn" } } }`
  data3 := `{"item_data": { "library": { "value": "Department", "desc": "UO Departmental Library" }, "location": { "value": "zartmus", "desc": "Music" } } }`
  path1 := "/almaws/v1/bibs/123/holdings/456/items/7890"
  path3 := "/almaws/v1/bibs/123/holdings/456/items/7892"
  //testserver responds to request from CheckLibrary calls
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path1 {
      fmt.Fprintf(w, data1)
    } else if r.URL.Path == path3 {
      fmt.Fprintf(w, data3)
    }
  }))
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")
  link1 := ts.URL + path1
  link3 := ts.URL + path3
  e1 := Eligible{}
  result,err := EligibleToUnlinkSuppressUnset([]string{link1, link1}, e1)
  if err != nil { log.Println(err) }
  if result.Unlink != true {t.Errorf("example1, incorrect unlink")}
  if result.Suppress != true {t.Errorf("example1 incorrect suppress")}
  if result.Unset != true {t.Errorf("example1 incorrect unset")}

  e2 := Eligible{}
  result,err = EligibleToUnlinkSuppressUnset([]string{link3, link1}, e2)
  if err != nil { log.Println(err) }
  if result.Unlink != true {t.Errorf("example3 incorrect unlink")}
  if result.Suppress != false {t.Errorf("example3 incorrect suppress")}
  if result.Unset != true {t.Errorf("example 3 incorrect unset")}
}

func TestHandleSerial(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  src, err := os.Open(homedir + "/fixtures/response_1743708271877.json")
  if err != nil { t.Errorf(err.Error()) }
  data,_ := io.ReadAll(src)

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //respond to request for holdings
    fmt.Fprintf(w, string(data))
  }))

  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")
  e := Eligible{Locations: []string{"kshort"}}
  e2, _ := HandleSerial("99126837001852", e)
  if !e2.SerialRequiresAction { t.Errorf("incorrect setting of serial flag") }
}


func TestHandleCases(t *testing.T){
  path1 := "/bibs/99126837001852"
  path2 := "/bibs/99126837001852/holdings"
  homedir := os.Getenv("HOME_DIR")
  src, err := os.Open(homedir + "/fixtures/fakeBW1_1743689261042.json")
  if err != nil { t.Errorf(err.Error()) }
  data1,_ := io.ReadAll(src)

  src, err = os.Open(homedir + "/fixtures/response_1743708271877.json")
  if err != nil { t.Errorf(err.Error()) }
  data2,_ := io.ReadAll(src)

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path1 {
      fmt.Fprintf(w, string(data1))
    } else if r.URL.Path == path2 {
      fmt.Fprintf(w, string(data2))
    }
  }))
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")

  e := Eligible{ Locations: []string{"kshort"}, Unlink: true }
  e2, _ := HandleCases("99126837001852", e)
  if e2.SerialRequiresAction != true { t.Errorf("serial flag not set correctly") }
  if e2.BoundWith != true { t.Errorf("boundwith flag not set correctly") }
  if e2.BoundWithMult != "123123123123" { t.Errorf("boundwith list not set correctly") }
  if e2.Unlink != false { t.Errorf("unlink is not set correctly ") }
  if len(e2.Locations) != 1 { t.Errorf("locations not preserved") }
}
