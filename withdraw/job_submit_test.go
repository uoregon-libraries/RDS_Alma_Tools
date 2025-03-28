package withdraw

import(
  "testing"
  "net/http/httptest"
  "net/http"
  "os"
  "fmt"
  "io"
  "io/ioutil"
  "strings"
  "rds_alma_tools/file"
)

func TestCheckJob(t *testing.T){
  os.Setenv("JOB_WAIT_TIME", "2s")
  os.Setenv("JOB_MAX_TRIES", "3")

  jobpath1 := "/almaws/v1/conf/jobs/112358/instances/98765"
  jobpath2 := "/almaws/v1/conf/jobs/112358/instances/98766"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == jobpath1 { 
      fmt.Fprintf(w, `{"status":{"value":"COMPLETED_SUCCESS"}}`)
    } else if r.URL.Path == jobpath2 {
      fmt.Fprintf(w, `{"status":{"value":"QUEUED"}, "alert":{"value": "BUSY"}}`)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  filename := "testcheckjob"
  list := map[string]Eligible{"112358":Eligible{Unlink: true, Suppress: true, Oclc: "1234"}}

  joblink1 := ts.URL + jobpath1
  CheckJob(joblink1, DummyJob, filename, list)
  joblink2 := ts.URL + jobpath2
  CheckJob(joblink2, nil, filename, list)
  report_dir := os.Getenv("REPORT_DIR")
  filepath := report_dir + "/" + filename
  content, err := ioutil.ReadFile(filepath)
  if err != nil { t.Errorf("unable to read report") }
  if !strings.Contains(string(content), "From dummy job") { t.Errorf("report is incorrect") }
  if !strings.Contains(string(content), "BUSY") { t.Errorf("report is incorrect") }
  _ = os.Remove(filepath)
}

func DummyJob(filename string, list map[string]Eligible){
  str := "From dummy job: "
  for k,_ := range list{ str += k }
  file.WriteReport(filename, []string{str})
}

func TestSubmitJob(t *testing.T){
  data := `{"parameter":[{"name":{"value":"set_id"},"value":"56789"}]}`
  return_data := `{"additional_info":{"link":"http://example.org/12121212"}}`
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    if len(body) != len(data) { t.Errorf("body of request wrong") }
    fmt.Fprintf(w, return_data)
  }))
  defer ts.Close()

  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  os.Setenv("ALMA_KEY", "123123")
  os.Setenv("TEST_URL", ts.URL)

  var params = []Param{
    Param{ Name: Val{ Value: "set_id" }, Value: "56789" },
  }
  link,_ := SubmitJob("123456", params)
  if link != "http://example.org/12121212" { t.Errorf("link was not returned") }
}

