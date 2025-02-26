package connect

import (
  "testing"
  "fmt"
  "os"
  "errors"
)

func TestResponses(t *testing.T){
  var r Report

  r.Responses = append(r.Responses, Response{ "http://xy.org/abc", BuildErrorMessage("Too many things")})

  r.Responses = append(r.Responses, Response{ "http://xy.org/def", ExtractMessage(`{"status":"success"}`)})

  r.Responses = append(r.Responses, Response{ "http://xy.org/ghi", BuildMessage("process completed")})

  r.Responses = append(r.Responses, Response{ "http://xy.org/jkl", ExtractAlmaError(`{"errorList":{"error":[{"errorMessage": "Input is not valid"}]}}`)})
  fmt.Println(r.ResponsesToString())
  // todo: parse the string
}

func TestWriteReport(t *testing.T){
  var r Report
  r.Responses = append(r.Responses, Response{ "http://xy.org/abc", BuildErrorMessage("Too many things")})
  filename := "banana"
  r.WriteReport(filename)
  path := os.Getenv("REPORT_DIR") + "/" + filename
  _, err := os.Stat(path)
  if err != nil { 
    if errors.Is(err, os.ErrNotExist) { t.Errorf("did not write report") }
  } else {
    _ = os.Remove(path)
  }
}
