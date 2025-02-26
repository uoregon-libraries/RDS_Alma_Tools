package connect

import (
  "log"
  "github.com/tidwall/gjson"
  "encoding/json"
  "os"
  "fmt"
)

type Report struct {
  Responses []Response
}

type Response struct {
  Id string
  Message Message
}

type Message map[string]any

func (r Response) ResponseToString() string{
  var output []byte
  var err error
  output, err = json.Marshal(r.Message)
  if err != nil { log.Println(err); return `{"id":` + r.Id + `", "error": "unable to marshal message" }` }
  return `{"id":"` + r.Id + `", "report":` + string(output) + "}"
}

func (r Report) ResponsesToString() string {
  all_resp := ""
  for _, elt := range r.Responses {
    all_resp += elt.ResponseToString()
  }
  return all_resp
}

func (r Report) WriteReport(filename string) {
  //writes to directory
  dir := os.Getenv("REPORT_DIR")
  path := fmt.Sprintf("%s/%s", dir, filename)
  f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  if err != nil { log.Println(err) }
  _,err = fmt.Fprintln(f, r.ResponsesToString())
  if err != nil { log.Println(err) }
}
// for reporting internal string messages
func BuildMessage(message string) Message{
  var m Message
  e := `{"message":"` + message + `"}`
  _ = json.Unmarshal([]byte(e), &m)
  return m

}
// for reporting results that are json
func ExtractMessage(message string) Message{
  var m Message
  err := json.Unmarshal([]byte(message), &m)
  if err != nil { log.Println(err); return BuildErrorMessage("unable to unmarshal message") }
  return m
}

//for reporting internal errs
func BuildErrorMessage(message string) Message{
  var m Message
  e := `{"error":"` + message + `"}`
  _ = json.Unmarshal([]byte(e), &m)
  return m
}

// for reporting Alma errs
func ExtractAlmaError(message string) Message{
  var m Message
  errMesses := gjson.Get(message, "errorList.error.#.errorMessage")
  for _, e := range errMesses.Array(){
    mess := `{"error":"` + e.String() + `"}`
  _ = json.Unmarshal([]byte(mess), &m)
  }
  return m
}
