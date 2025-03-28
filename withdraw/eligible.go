package withdraw

import(
  "rds_alma_tools/connect"
  "github.com/tidwall/gjson"
  "bytes"
)

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

func UniqueBibs(data []byte)map[string]bool{
  unique := map[string]bool{}
  lines := bytes.Split(data, []byte("\n"))
  for _, line := range lines{
    if string(line) == "" { break }
    lineMap := LineMap(string(line))
    unique[lineMap["mms_id"]] = true
  }
  return unique
}

// runs logic for one item
func CheckLibrary(link string)([]bool, error){
  params := []string{"", ApiKey()}
  item, err := connect.Get(link, params)
  if err != nil { return nil, err }
  library := gjson.GetBytes([]byte(item), "item_data.library.value")
  if library.String() == "Withdrawn" {
    return []bool{ true, true }, nil } else if library.String() == "Department" {
    return []bool{ true, false }, nil } else {
    return []bool{ false, false }, nil
  }
}

// returns one result based on checking all items 
func EligibleToUnlinkAndSuppress(items []string)([]bool, error){
  suppress := true
  for _, v:= range items{
    arr, err := CheckLibrary(v)
    if err != nil { return nil, err }
    if arr[0] != true { return []bool{false,false}, nil }
    if arr[1] != true { suppress = false }
  }
  return []bool{true, suppress}, nil
}

//generates list of bibs to unlink or unlink and suppress
//returns map of mmsid keys and []bool
func EligibleToUnlinkAndSuppressList(data []byte)(map[string][]bool, error){
  var eligibleList = map[string][]bool{}
  bibs := UniqueBibs(data)
  for k,_ := range bibs{
    items, err := BibItems(k)
    if err != nil { return nil, err }
    arr, err := EligibleToUnlinkAndSuppress(items)
    if err != nil { return nil, err }
    if arr[0] { eligibleList[k] = arr }
  }
  return eligibleList, nil
}
