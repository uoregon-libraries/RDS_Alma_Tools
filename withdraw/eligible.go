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
  SerialRequiresAction bool // requires further handling after withdraw process
  Locations []string
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

func HandleCases(mmsId string, eligible Eligible)(Eligible, error){
  url := BuildBibLink(mmsId)
  params := []string{ ApiKey() }
  json, err := connect.Get(url, params)
  if err != nil { return eligible, err }
  serial := Is_serial(json)
  if serial {
    e, err := Handle_serial(mmsId, eligible)
    return e, err
  }
  return eligible, nil
}

func Handle_serial(mmsId string, eligible Eligible)(Eligible, error){
  holding_json, err := Holdings(mmsId)
  if err != nil { return eligible, err }
  eligible.SerialRequiresAction = false
  tr := gjson.GetBytes(holding_json, "total_record_count")
  if tr.Int() == 1 { return eligible, nil }
  holdings := gjson.GetBytes(holding_json, "holding")
  for _, h := range holdings.Array() {
    h_loc := gjson.Get(h.String(), "location.value")
    for _, l := range eligible.Locations {
      if h_loc.String() == l { eligible.SerialRequiresAction = true }
    }
  }
  return eligible, nil
}

func Holdings(mmsId string)([]byte,error){
  url := BuildHoldingLink(mmsId, "")
  params := []string{ ApiKey() }
  json, err := connect.Get(url, params)
  if err != nil { return nil, err }
  return json, nil
}

func UniqueBibs(data []byte) map[string]Eligible {
  unique := map[string]Eligible{}
  lines := bytes.Split(data, []byte("\n"))
  for _, line := range lines{
    if string(line) == "" { break }
    lineMap := LineMap(string(line))
    e, ok := unique[lineMap["mms_id"]]
    if ok {
      e.Locations = append(e.Locations, lineMap["location"]) } else {
      e = Eligible{ Oclc: lineMap["oclc"], Locations: []string{ lineMap["location"] } }
    }
    unique[lineMap["mms_id"]] = e
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
    eligible.Oclc = v.Oclc
    eligible, err = HandleCases(k, eligible)
    if err != nil { errs =  append(errs, fmt.Sprintf("Eligibility error: %s", k)); continue }
    eligibleList[k] = eligible
  }
  return eligibleList, errs
}
