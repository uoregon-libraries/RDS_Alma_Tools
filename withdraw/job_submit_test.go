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
      fmt.Fprintf(w, `{"status":{"value":"QUEUED"}}`)
    }
  }))
  defer ts.Close()
  os.Setenv("ALMA_URL", ts.URL + "/almaws/v1/")
  os.Setenv("VERBOSE", "true")
  os.Setenv("DEBUG", "true")
  filename := "testcheckjob"
  list := map[string][]bool{"112358":[]bool{true,true}}

  joblink1 := ts.URL + jobpath1
  CheckJob(joblink1, DummyJob, filename, list)
  joblink2 := ts.URL + jobpath2
  CheckJob(joblink2, nil, filename, list)
  report_dir := os.Getenv("REPORT_DIR")
  filepath := report_dir + "/" + filename
  content, err := ioutil.ReadFile(filepath)
  if err != nil { t.Errorf("unable to read report") }
  if !strings.Contains(string(content), "From dummy job") { t.Errorf("report is incorrect") }
  if !strings.Contains(string(content), "Unable to confirm") { t.Errorf("report is incorrect") }
  _ = os.Remove(filepath)
}

func DummyJob(filename string, list map[string][]bool){
  str := "From dummy job: "
  for k,_ := range list{ str += k }
  WriteReport(filename, str)
}

func TestSubmitJob(t *testing.T){
  data := `{"Parameter":[{"Name":{"Value":""},"Value":""}]}`
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
  filename := "testsubmitjob"
  link,_ := SubmitJob(filename, "123456")
  if link != "http://example.org/12121212" { t.Errorf("link was not returned") }
}
