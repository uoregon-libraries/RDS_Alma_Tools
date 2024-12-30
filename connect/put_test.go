package connect

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "io"
)

func TestPut(t *testing.T){
  homedir := os.Getenv("HOME_DIR")
  data, err := os.ReadFile(homedir + "/fixtures/item_record.json")
  if err != nil { t.Errorf(err.Error())}
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/almaws/v1/banana" { t.Errorf("wrong path") }
    body, _ := io.ReadAll(r.Body)
    if len(body) != len(data) { t.Errorf("body of request wrong") }
    if v := r.URL.Query().Get("apikey"); v != "123456" { t.Errorf("params wrong") }
    fmt.Fprintf(w, "hello hello")
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("VERBOSE", "true")
  final_url := ts.URL + "/almaws/v1/banana"
  os.Setenv("DEBUG", "true")
  os.Setenv("TEST_URL", final_url + "?apikey=123456")
  params := []string{ "apikey=123456" }
  resp, err := Put(final_url, params, string(data))
  if err != nil { t.Errorf(err.Error()) }
  if string(resp) != "hello hello" { t.Errorf("wrong response") }
}
