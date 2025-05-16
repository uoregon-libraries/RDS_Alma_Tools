package withdraw

import(
  "github.com/tidwall/gjson"
  "github.com/beevik/etree"
  "fmt"
  "strings"
)

func ExtractJsonField(field string, json_data []byte)string{
  return gjson.GetBytes(json_data, field).String()
}

func Leader(marc_string string)string{
  marc_tree := MarcTree(marc_string)
  leader := marc_tree.FindElement("//leader")
  if leader == nil { return "" }
  return leader.Text()
}

func Subfield(tag string, subfield string, marc_string string)string{
  marc_tree := MarcTree(marc_string)
  sub_elt := marc_tree.FindElement(fmt.Sprintf("//datafield[@tag='%s']/subfield[@code='%s']", tag, subfield))
  if sub_elt == nil { return "" }
  return sub_elt.Text()
}

func Datafield(tag string, marc_string string)string{
  marc_tree := MarcTree(marc_string)
  tag_elt := marc_tree.FindElement(fmt.Sprintf("//datafield[@tag='%s']", tag))
  if tag_elt == nil { return "" }
  return tag_elt.Text()
}

func Controlfield(tag string, marc_string string)string{
  marc_tree := MarcTree(marc_string)
  tag_elt := marc_tree.FindElement(fmt.Sprintf("//controlfield[@tag='%s']", tag))
  if tag_elt == nil { return "" }
  return tag_elt.Text()

}
func IsSerial(json_data []byte)bool{
  marc := ExtractJsonField("anies.0", json_data)
  leader := []byte(Leader(marc))
  if leader[7] == 's' || leader[7] == 'i' { return true }
  return false
}

// returns false, empty string if not bound with
// returns true, string of bibs if bound with multiple
// returns true, empty string if bound with not multiple
func IsBoundWith(json_data []byte)(bool, string){
  marc := ExtractJsonField("anies.0", json_data)
  tag962 := Subfield("962", "a", marc)
  if tag962 != "BoundWithRecord" { return false, "" }

  title := Subfield("222", "a", marc)
  tag001 := Controlfield("001", marc)
  if ((strings.HasPrefix(title, "Multiple")) && (tag001 == "")) {
    tag774 := Subfield("774", "w", marc)
    return true, tag774
  }
  return true, ""
}

func MarcTree(marc_string string)*etree.Document{
  marc_tree := etree.NewDocument()
  marc_tree.ReadFromString(marc_string)
  return marc_tree
}
