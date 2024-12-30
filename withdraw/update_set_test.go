package withdraw

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "io"
)

func TestUpdateSet(t *testing.T){
  data := `{"Name":"banana","Type":{"Value":"LOGICAL"},"Content":{"Value":"BIB_MMS"},"Query":{"Value":""},"Members":{"Member":[{"Id":"112358"}]}}`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    if len(body) != len(data) { t.Errorf("body of request wrong") }
    fmt.Fprintf(w, "hello hello")
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  setname := "banana"
  setid := "12345"
  os.Setenv(setname, setid)
  final_url := ts.URL + "/almaws/v1/conf/sets/" + setid

  os.Setenv("TEST_URL", final_url + "?apikey=123456&op=replace_members")
  filename := "testupdateset"
  err := UpdateSet(filename, "banana", "BIB_MMS", []string{"112358"})
  if err != nil { t.Errorf(err.Error()) }
}
