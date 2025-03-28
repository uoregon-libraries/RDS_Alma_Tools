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

func TestCheckLibrary(t *testing.T){
  data1 := `{"item_data": { "library": { "value": "Withdrawn", "desc": "Withdrawn Library" } } }`
  data2 := `{"item_data": { "library": { "value": "Science", "desc": "Price Science Commons" } } }`
  data3 := `{"item_data": { "library": { "value": "Department", "desc": "UO Departmental Library" } } }`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/bibs/123456789" {
      fmt.Fprintf(w, data1)
    } else if r.URL.Path == "/bibs/123456788" {
      fmt.Fprintf(w, data2)
    } else {
      fmt.Fprintf(w, data3)
    }
   }))
  defer ts.Close()

  result, _ := CheckLibrary(ts.URL + "/bibs/123456789")
  if !slices.Equal(result, []bool{true,true}) {t.Errorf("should be true/true")}

  result, _ = CheckLibrary(ts.URL + "/bibs/123456788")
  if !slices.Equal(result, []bool{false,false}) {t.Errorf("should be false/false")}

  result, _ = CheckLibrary(ts.URL + "/bibs/123456787")
  if !slices.Equal(result, []bool{true,false}) {t.Errorf("should be true/false")}
}

func TestUniqueBibs(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  src, err := os.Open(homedir + "/fixtures/export.tsv")
  data,_ := io.ReadAll(src)
  if err != nil { t.Fatalf("did not read file") }
  bibs := UniqueBibs(data)
  if len(bibs) != 2 { t.Errorf("should be size 2") }
  fmt.Println(bibs)
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

func TestEligibleToUnlinkAndSuppress(t *testing.T){
  data1 := `{"item_data": {"library": {"value":"Withdrawn"}}}`
  data2 := `{"item_data": {"library": {"value":"Department"}}}`
  data3 := `{"item_data": {"library": {"value":"Banana"}}}`
  path1 := "/almaws/v1/bibs/123/holdings/456/items/7890"
  path2 := "/almaws/v1/bibs/123/holdings/456/items/7891"
  path3 := "/almaws/v1/bibs/123/holdings/456/items/7892"
  //testserver responds to request from CheckLibrary calls
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == path1 {
      fmt.Fprintf(w, data1)
    } else if r.URL.Path == path2 {
      fmt.Fprintf(w, data2)
    } else if r.URL.Path == path3 {
      fmt.Fprintf(w, data3)
    }
  }))
  os.Setenv("ALMA_URL", ts.URL + "/")
  os.Setenv("ALMA_KEY", "almakey")
  link1 := ts.URL + path1
  link2 := ts.URL + path2
  link3 := ts.URL + path3
  result,_ := EligibleToUnlinkAndSuppress([]string{link1, link1})
  if !slices.Equal(result, []bool{true, true}){t.Errorf("incorrect result")}
  result,_ = EligibleToUnlinkAndSuppress([]string{link1, link2})
  if !slices.Equal(result, []bool{true, false}){t.Errorf("incorrect result")}
  result,_ = EligibleToUnlinkAndSuppress([]string{link1, link3})
  if !slices.Equal(result, []bool{false, false}){t.Errorf("incorrect result")}
}
