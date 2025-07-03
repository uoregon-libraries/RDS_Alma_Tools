package oclc

import(
  "fmt"
  "rds_alma_tools/file"
  "log"
  "strconv"
  "encoding/json"
)

func UnsetHoldings(filename string, list map[string]string){
  token,err := OclcAuth()
  if err != nil { log.Println(err); file.WriteReport(filename, []string{ "unable to authenticate with OCLC" }); return }
  results := []string{}
  for k,v := range list{
    resp, err := UnsetHolding(v, token)
    if err != nil { 
      if resp != "" {
        results = append(results, k + ": " + resp) } else {
        results = append(results, k + ": " + err.Error()) }
    }
    results = append(results, k + ": " + "oclc unset success")
  }
  file.WriteReport(filename, results)
}

func UnsetHolding(oclc_num string, token string)(string, error){
  url := fmt.Sprintf("manage/institution/holdings/%s/unset", oclc_num)
  oclc_resp, err := Request(token, "POST", "", url, "", "json")
  return oclc_resp, err
}

func CheckHolding(oclc_num string)(string, error){
  token,err := OclcAuth()
  if err != nil { log.Println(err); return "", err }

  url := fmt.Sprintf("manage/institution/holdings/current?oclcNumbers=%s", oclc_num)
  oclc_resp, err := Request(token, "GET", "", url, "", "json")
  if err != nil { return "", err }
  var h Holdings
  err = json.Unmarshal([]byte(oclc_resp), &h)
  return strconv.FormatBool(h.Holding_set), err
}

type Holdings struct{
  Holding_set bool `json:"holdingSet"`
}
