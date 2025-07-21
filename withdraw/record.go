package withdraw

import (
  "fmt"
  "strings"
  "regexp"
)

type Record struct{
  Bib_data Bib `json:"bib_data"`
  Holding_data Holding `json:"holding_data"`
  Item_data Item `json:"item_data"`
}

func (r Record)Stringify()string{
  str := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s",
    r.Bib_data.Mms_id, 
    r.Bib_data.Title,
    OclcSelect(r.Bib_data.Network),
    r.Holding_data.Holding_id,
    r.Item_data.Item_pid,
    r.Item_data.Barcode,
    r.Item_data.Base_status.Desc,
    r.Item_data.Library.Value,
    r.Item_data.Location.Value,
    r.Item_data.Public_note,
    r.Item_data.Fulfillment_note,
    r.Item_data.Internal_note_1,
    r.Item_data.Internal_note_2,
    r.Item_data.Internal_note_3,
    r.Item_data.Retention_note)
  return str
}

type Bib struct{
  Mms_id string `json:"mms_id"`
  Title string `json:"title"`
  Bib_suppress string `json:"bib_suppress_from_publishing"`
  Network []string `json:"network_number"`
}

type Holding struct{
  Holding_id string `json:"holding_id"`
}

type Item struct{
  Item_pid string `json:"pid"`
  Barcode string `json:"barcode"`
  Base_status Desc `json:"base_status"`
  Library Value `json:"library"`
  Location Value `json:"location"`
  Public_note string `json:"public_note"`
  Fulfillment_note string `json:"fulfillment_note"`
  Internal_note_1 string `json:"internal_note_1"`
  Internal_note_2 string `json:"internal_note_2"`
  Internal_note_3 string `json:"internal_note_3"`
  Retention_note string `json:"retention_note"`
}

type Value struct{
  Value string `json:"value"`

}
type Desc struct{
  Desc string `json:"desc"`
}

func LineMap(line string)map[string]string{
  m := map[string]string{}
  keys := strings.Split(BriefItemHead(), "\t")
  arr := strings.Split(line, "\t")
  for ind, key := range keys{
    m[key] = arr[ind]
  }
  return m
}

func AllianceLinked(vals []string)string{
  re := regexp.MustCompile(`\(EXLNZ-01ALLIANCE_NETWORK\)`)
  for _, val := range vals {
    arr := re.FindStringSubmatch(val)
    if arr != nil {
      return "true"
    }
  }
  return "false"
}

func OclcSelect(vals []string)string{
  re := regexp.MustCompile(`\(OCoLC\)[a-z]*([0-9]*)`)
  for _, val := range vals {
    arr := re.FindStringSubmatch(val)
    if arr != nil {
      return arr[1]
    }
  }
  return ""
}
