package connect

import (

  "os"
  "net/http"
  "log"
  "time"
  "errors"
  "io"
  "strings"
)

//url /almaws/v1/bibs/<mms_id>/holdings/<holding_id>/items/<item_id>
//params: apikey=abcde12341234
func Put(_url string, params []string, json_record string)([]byte, error){
  body, err := push("PUT", _url, params, json_record)
  return body, err
}

func Post(_url string, params []string, json_record string)([]byte, error){
  body, err := push("POST", _url, params, json_record)
  return body, err
}

func push(method string, _url string, params []string, json_record string)([]byte, error){
  verbose := os.Getenv("VERBOSE")
  debug := os.Getenv("DEBUG")
  param_str := strings.Join(params[:], "&")
  final_url := _url + "?" + param_str
  if debug == "true" {
    log.Println("Swapping " + final_url + "for test url")
    final_url = os.Getenv("TEST_URL")
  }
  data := strings.NewReader(json_record)
  req, err := http.NewRequest(method, final_url, data)
  if err != nil { log.Println(err); return nil, errors.New("unable to create http request") }
  req.Header.Set("accept", "application/json")
  RequestDump(verbose, req)
  client := &http.Client{
    Timeout: time.Second * 60,
  }
  response, err := client.Do(req)
  ResponseDump(verbose, response)
  defer response.Body.Close()
  if err != nil { log.Println(err); return nil, errors.New("unable to complete http request") }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return nil, errors.New("unable to read response from alma") }
  if response.StatusCode != 200 { return body, errors.New("alma errors") }

  return body, nil
}
