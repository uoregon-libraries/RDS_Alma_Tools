package withdraw

import(
  "rds_alma_tools/connect"
  "github.com/tidwall/gjson"
  "bytes"
  "os"
  "io"
  "fmt"
  "errors"
)

type Eligible struct {
  Unlink   bool
  Suppress bool
  Unset    bool
  Oclc     string
}

//The only errors will be from connect.Get
//propagate and return nil

// retrieves all item links for a given bib
func BibItems(mmsId string)([]string, error) {
  arr := []string{}
  url := BuildItemLink(mmsId, "all", "")
  params := []string{ ApiKey() }
  json, err := connect.Get(url, params)
  if err != nil { return nil, err }
  members := gjson.GetBytes(json, "item.#.link")
  for _, link := range members.Array(){
    arr = append(arr, link.String())
  }
  return arr, nil

}

func UniqueBibs(data []byte)map[string]string{
  unique := map[string]string{}
  lines := bytes.Split(data, []byte("\n"))
  for _, line := range lines{
    if string(line) == "" { break }
    lineMap := LineMap(string(line))
    unique[lineMap["mms_id"]] = lineMap["oclc"]
  }
  return unique
}

type locationVals struct{
  NZ string
  Primo string
  ORU string
  UOL string
}

type LLKey struct{
  LibCode string
  LocCode string
}

func LibraryLocationMap()map[LLKey]locationVals{
  locmap := map[LLKey]locationVals{}
  homedir := os.Getenv("HOME_DIR")
  src,_ := os.Open(homedir + "/withdraw/location_eligibility.txt")
  data,_ := io.ReadAll(src)
  lines := bytes.Split(data, []byte("\n"))
  for _,line := range lines{
    if string(line) == "" { break }
    arr := bytes.Split(line, []byte("\t"))
    locmap[LLKey{string(arr[0]), string(arr[1])}] = locationVals{NZ: string(arr[2]), Primo: string(arr[3]), ORU: string(arr[4]), UOL: string(arr[5])}
  }
  return locmap
}

func ItemLibraryLocation(link string)(LLKey, error){
  params := []string{"", ApiKey()}
  item, err := connect.Get(link, params)
  if err != nil { return LLKey{}, err }
  library := gjson.GetBytes([]byte(item), "item_data.library.value")
  location := gjson.GetBytes([]byte(item), "item_data.location.value")
  return LLKey{LibCode: library.String(), LocCode: location.String()}, nil
}

func EligibleToUnlinkSuppressUnset(items []string)(Eligible, error){
  locmap := LibraryLocationMap()
  e := Eligible{Unlink: true, Suppress: true, Unset: true}
  for _, v:= range items{
    k,err := ItemLibraryLocation(v)
    if err != nil { return Eligible{}, err }
    chart := locmap[k]
    if chart.NZ == "" { return Eligible{}, errors.New("Eligibility not known") }
    if chart.NZ == "Y" { e.Unlink = false }
    if chart.Primo == "Y" { e.Suppress = false }
    if chart.ORU == "Y" { e.Unset = false }
  }
  return e, nil
}

func EligibleToUnlinkSuppressUnsetList(data []byte)(map[string]Eligible, []string){
  var eligibleList = map[string]Eligible{}
  bibs := UniqueBibs(data)
  errs := []string{}
  for k,v := range bibs{
    items, err := BibItems(k)
    if err != nil { errs = append(errs, fmt.Sprintf("Eligibility error: %s", k)); continue }
    eligible, err := EligibleToUnlinkSuppressUnset(items)
    if err != nil { errs = append(errs, fmt.Sprintf("Eligibility error: %s", k)); continue }
    eligible.Oclc = v
    eligibleList[k] = eligible
  }
  return eligibleList, errs
}
