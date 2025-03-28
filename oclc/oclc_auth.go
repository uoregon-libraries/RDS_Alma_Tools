package oclc

import(
  "net/http"
  "encoding/json"
  "os"
  "log"
  "time"
  "io/ioutil"
  "github.com/labstack/echo/v4"
  "rds_alma_tools/utils"
  "net/http/httputil"
  "errors"
)

type OclcToken struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

func OclcAuth() (string, error) {
  url := os.Getenv("OCLC_AUTH_URL")
  name := os.Getenv("OCLC_NAME")
  pass := os.Getenv("OCLC_PASS")
  verbose := os.Getenv("VERBOSE")

  req, err := http.NewRequest("POST", url, nil)
  if err != nil { log.Println(err); return "", errors.New("unable to create request") }
  req.SetBasicAuth(name, pass)
  if verbose == "true" {
    reqdump, err := httputil.DumpRequestOut(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }

  client := &http.Client{
    Timeout: time.Second * 10,
  }
  response, err := client.Do(req)
  if err != nil { log.Println(err); return "", errors.New("unable to complete http request") }
  defer response.Body.Close()
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }

  byteVal, _ := ioutil.ReadAll(response.Body)
  var ot OclcToken
  err = json.Unmarshal(byteVal, &ot)
  if err != nil { log.Println(err); return "", errors.New("unable to extract token") }
  return ot.AccessToken, nil
}


func GetToken(c echo.Context) (string, error){
  token, err := utils.FetchCookieVal(c, "oclc_token")
  if token == "" || err != nil {
    token, err = OclcAuth()
    if err != nil { log.Println(err); return "", err }
    utils.WriteCookie(c, 20, "oclc_token", token)
  }
  return token, nil
}


