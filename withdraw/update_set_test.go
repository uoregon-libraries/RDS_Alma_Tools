package withdraw

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "io"
  "slices"
)

func TestUpdateSet(t *testing.T){
  data := `{"type":{"value":"ITEMIZED"},"content":{"value":"BIB_MMS"},"members":{"member":[{"id":"112358"}]}}`
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

  eligibleList := map[string]Eligible{"112358": Eligible{Suppress: true, Unlink: true, Oclc: "54321"}}
  err := UpdateSet("banana", "BIB_MMS", eligibleList)
  if err != nil { t.Errorf(err.Error()) }
}

func TestSetMembers( t *testing.T){
  eligiblelist := map[string]Eligible{
    "a": Eligible{Suppress: true, Unlink: true, Oclc: "54321"},
    "b": Eligible{Suppress: true, Unlink: true, Oclc: "54322"},
    "c": Eligible{Suppress: true, Unlink: true, Oclc: "54323"}}
  set := InitSet("BIB_MMS")
  set = SetMembers(set, eligiblelist)
  compare := []string{"a", "b", "c"}
  for _,v := range set.Members.Member{
    if slices.Contains(compare, v.Id) != true { t.Errorf("incorrect set membership") }
  }
}
