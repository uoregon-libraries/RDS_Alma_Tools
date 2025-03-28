package oclc

import(
  "testing"
  "strings"
  //"regexp"
  "net/http/httptest"
  "net/http"
  "os"
)

func TestAssembleUrl(t *testing.T){
  parts_no_id := []string{"https://blah.org","path",""}
  parts_id := []string{"https://blah.org","path","abcd"}
  url_no_id := assembleUrl(parts_no_id)
  url_id := assembleUrl(parts_id)
  if url_no_id != "https://blah.org/path" {t.Fatalf("assembled url is incorrect")}
  if url_id != "https://blah.org/path/abcd" {t.Fatalf("assembled url is incorrect")}
}

func TestRequestUpdate(t *testing.T){
   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { t.Fatalf("request action is not correct") }
    if r.Header.Get("accept") != "application/json" { t.Fatalf("request header is not correct")}
    if r.Header.Get("Content-Type") != "" { t.Fatalf("request header is not correct")}
    arr := strings.Split(r.URL.String(), "/")
    pathend := arr[len(arr)-1]
    if pathend != "unset" { t.Fatalf("url is incorrect") }
  }))
  defer ts.Close()
  os.Setenv("OCLC_URL", ts.URL + "/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  _,_ = Request("token", "POST", "", "manage/institution/holdings/123/unset", "", "json")
}
