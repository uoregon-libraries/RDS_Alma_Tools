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

func Put(url string, params []string, json_record string)([]byte, error){
  verbose := os.Getenv("VERBOSE")
  param_str := strings.Join(params[:], "&")
  final_url := url + "?" + param_str
  data := strings.NewReader(json_record)
  req, err := http.NewRequest("PUT", final_url, data)
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
  if response.StatusCode != 200 { return nil, errors.New(string(body)) }

  return body, nil
}
