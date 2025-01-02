package withdraw

import(
  "github.com/tidwall/gjson"
  "time"
  "rds_alma_tools/connect"
  "encoding/json"
  "log"
  "fmt"
  "strings"
  "os"
  "strconv"
)

// filename, list
type ProcessFunc func(string, map[string][]bool)

func CheckJob(joblink string, nextFun ProcessFunc, filename string, eligibleList map[string][]bool ){
  MAX, _ := strconv.Atoi(os.Getenv("JOB_MAX_TRIES"))
  span,_ := time.ParseDuration(os.Getenv("JOB_WAIT_TIME"))
  i := 0
  params := []string{ ApiKey() }

  for i < MAX {
    resp,err := connect.Get(joblink, params)
    if err != nil { 
      log.Println(err)
      /*** count this as one try ****/
      i += 1
      time.Sleep(span)
      continue
    }
    result := ExtractJobResults(resp)
    if result == "false"{
      i += 1
      time.Sleep(span)
      continue
    } else if result == "COMPLETED_SUCCESS" {
      if nextFun != nil {
        nextFun(filename, eligibleList)
      }
    }
    WriteReport(filename, result)
    return
  }
  WriteReport(filename, "Unable to confirm that job completed: " + joblink)
}

func ExtractJobResults(resp []byte)string{
  //the docs on using progress are unclear
  status := gjson.GetBytes(resp, "status.value")
  if !strings.Contains(status.String(), "COMPLETED") { return "false" }
  if status.String() == "COMPLETED_SUCCESS" { return status.String() }
  alert := gjson.GetBytes(resp, "alert.value")
  return alert.String()
}

func SubmitJob(filename string, jobid string)(string, error){
  url := fmt.Sprintf("%s/jobs/%s", BaseUrl(), jobid)
  params := []string{ "op=run", ApiKey() }
  job := JobInit()
  json,_ := json.Marshal(job)
  resp,err := connect.Post(url, params, string(json))
  if err != nil { 
    log.Println(err)
    WriteReport(filename, err.Error())
    return "", err
  }
  link := ExtractJobInstance(resp)
  if err != nil { log.Println(err); WriteReport(filename, err.Error()); return "", err }
  return link, nil
}

func ExtractJobInstance(resp []byte)(string){
  instance := gjson.GetBytes(resp, "additional_info.link")
  return instance.String()
}

// boilerplate
// finish this when the actual jobs are available
func JobInit()Job{
  job := Job{ Parameter: []Param{Param{Name: Val{Value: ""}, Value: ""}}}
  return job
}

type Job struct{
  Parameter []Param
}

type Param struct{
  Name Val
  Value string
}

// Val declaration in update_set
