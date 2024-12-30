package connect

import (
  "testing"
  "fmt"
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
