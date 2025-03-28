package oclc

import(
  "fmt"
  "rds_alma_tools/file"
  "log"
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
