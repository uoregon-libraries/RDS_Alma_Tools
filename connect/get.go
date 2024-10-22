package connect

import(
  "url"
  "fmt"
  "os"
  "net/http"
  "log"
  "time"
  "errors"
  "strings"
)

//url /almaws/v1/conf/sets/setid
//params: api: "conf", param_ids: ["sets", <setid>]
//url /almaws/v1/bibs/<mms_id>/holdings/<holding_id>/items/<item_id>
//params: api: "bibs", param_ids: [<mms_id>, "holdings", <holding_id>, "items", <item_id>]

func Get(api string, param_ids []string)(string, error){
  verbose := os.Getenv("VERBOSE")
  base_url := os.Getenv("ALMA_URL")
  test := os.Getenv("TEST")
  url := url.JoinPath(base_url, api, strings.Join(param_ids[:],"/")) 

  req, err := http.NewRequest("GET", url)
  if err != nil { log.Println(err); return "", errors.New("unable to create http request") }
  req.Header.Set("accept", "application/json")

  RequestDump(verbose, req)
  client := &http.Client{
    Timeout: time.Second * 60,
  }

  if test == "true" { return fmt.Sprintf("{\"id\":\"%s\"}", id), nil }

  response, err := client.Do(req)
  ResponseDump(verbose, response)
  defer response.Body.Close()
  if err != nil { log.Println(err); return "", errors.New("unable to complete http request") }
  body, err := io.ReadAll(response.Body)
  if err != nil { log.Println(err); return "", errors.New("unable to read response from alma") }
  if response.StatusCode != 200 { return string(body), errors.New("alma errors") }

  return string(body), nil
}
