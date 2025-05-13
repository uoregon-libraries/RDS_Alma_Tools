package withdraw

import(
  "github.com/tidwall/gjson"
  "github.com/beevik/etree"
)

func Marc_parse(field string, json_data []byte)string{
  return gjson.GetBytes(json_data, field).String()
}

func Leader(marc_string string)string{
  marc_tree := etree.NewDocument()
  marc_tree.ReadFromString(marc_string)
  leader := marc_tree.FindElement("//leader")
  return leader.Text()
}

func Is_serial(json_data []byte)bool{
  marc := Marc_parse("anies", json_data)
  leader := []byte(Leader(marc))
  if leader[7] == 's' || leader[7] == 'i' { return true }
  return false
}
