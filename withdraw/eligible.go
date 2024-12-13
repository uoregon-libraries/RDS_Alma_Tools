package withdraw

import(
  "rds_alma_tools/connect"
  "github.com/tidwall/gjson"
  "bufio"
  "io"
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

func UniqueBibs(src io.Reader)map[string]bool{
  unique := map[string]bool{}
  scanner := bufio.NewScanner(src)
  for scanner.Scan(){
    line := scanner.Text()
    lineMap := LineMap(string(line))
    if unique[lineMap["mms_id"]] != true{
      unique[lineMap["mms_id"]] = true
    }
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

// runs logic for all items for a bib
func EligibleToUnlinkAndSuppress(mms_id string)([]bool, error){
  items, err := BibItems(mms_id)
  if err != nil { return nil, err }
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
func EligibleToUnlinkAndSuppressList(src io.Reader)(map[string][]bool, error){
  eligibleList := map[string][]bool{}
  bibs := UniqueBibs(src)
  for k,_ := range bibs{
    arr, err := EligibleToUnlinkAndSuppress(k)
    if err != nil { return nil, err }
    if arr[0] { eligibleList[k] = arr }
  }
  return eligibleList, nil
}
