package withdraw

import (

  "rds_alma_tools/connect"
  "github.com/tidwall/sjson"
  "os"
  "bufio"
  "github.com/tidwall/gjson"
  "strings"
  "io"
)

func UpdateItems(r connect.Report, loc_type string, src io.Reader)(connect.Report){
  //debug := os.Getenv("DEBUG")
  //process the data for updating item records
  scanner := bufio.NewScanner(src)
  for scanner.Scan(){
    line := scanner.Text()
    itemRec, response := UpdateItem(loc_type, line)
    if response.Id != ""{
      r.Responses = append(r.Responses, response)
      continue
    }
    params := []string{ ApiKey() }
    url := gjson.GetBytes(itemRec, "link").String()
    //if debug == "true" { continue }
    body, err := connect.Put(url, params, string(itemRec))
    if err != nil {
      if body != nil {
      r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.ExtractAlmaError(err.Error()) }) } else { 
      r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error()) }) }
    }
    if body != nil { r.Responses = append(r.Responses, connect.Response{ Id: url, Message: connect.BuildMessage("success") } )
    }
  }
  return r
}

func UpdateItem(newLocType string, data string)([]byte, connect.Response){
  lineMap := LineMap(data)
  url := BuildItemLink(lineMap["mms_id"], lineMap["holding_id"], lineMap["pid"])
  params := []string{ ApiKey() }
  itemRec, err := connect.Get(url, params)
  if err != nil { 
    if itemRec != nil { return nil, connect.Response{ Id: url, Message: connect.ExtractAlmaError(string(itemRec))} } else { return nil, connect.Response{ Id: url, Message: connect.BuildMessage(err.Error())} } }
  libMap := LoadMap()
  newLocVal := libMap[Key{lineMap["library"],newLocType,"value"}]
  newLocDesc := libMap[Key{lineMap["library"],newLocType,"desc"}]
  newLibVal := "Withdrawn" /******Logic?******/
  newLibDesc := "Withdrawn Library" /******************/
  internalNote3 := lineMap["internal_note_3"] + "|WD FY" + FiscalYear(TimeNow())
  //using sjson insert new library, location, append note
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.location.value", newLocVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.location.desc", newLocDesc)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.internal_note_3", internalNote3)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.value", newLibVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.desc", newLibDesc)
  return itemRec, connect.Response{Id:"", Message: connect.BuildMessage("")}
}

type Key struct{
  Libname, Libtype, Proptype string
}

func LoadMap()map[Key]string{
  libmap := map[Key]string{}
  homedir := os.Getenv("HOME_DIR")
  data,_ := os.Open(homedir + "/withdraw/library_map.txt")
  fileScanner := bufio.NewScanner(data)
  fileScanner.Split(bufio.ScanLines)
  for fileScanner.Scan(){
    arr := strings.Split(fileScanner.Text(), "\t")
    libmap[Key{arr[0], arr[1], arr[2]}] = arr[3]
  }
  return libmap
}
