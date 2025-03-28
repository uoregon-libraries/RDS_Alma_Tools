package withdraw

import (
  "encoding/json"
  "rds_alma_tools/connect"
  "os"
  "log"
  "net/url"
)

//setcontent either BIB_MMS or ITEM
func UpdateSet(setname string, setcontent string, eligibleList map[string]Eligible)error{
  setid := os.Getenv(setname)
  params := []string{ "op=replace_members", ApiKey() }
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("conf", "sets", setid)
  set := InitSet(setcontent)
  set = SetMembers(set, eligibleList)
  body, err := json.Marshal(set)
  if err != nil { log.Println(err); return err }
  _, err = connect.Post(_url.String(), params, string(body))
  if err != nil { log.Println(err); return err }
  return nil 
}

func SetMembers(set Set, eligibleList map[string]Eligible)Set{
  for k,_ := range eligibleList {
    set.Members.Member = append(set.Members.Member, RecId{ Id: k })
  }
  return set
}

func InitSet(content string) Set{
  var set = Set{Type: Val{Value: "ITEMIZED"}, Content: Val{Value: content}, Members: MemberArr{ Member: []RecId{} }}
  return set
}

type Set struct {
  Type Val          `json:"type"`
  Content Val       `json:"content"`
  Members MemberArr `json:"members"`
}

type Val struct {
  Value string `json:"value"`
}

type MemberArr struct {
  Member []RecId `json:"member"`
}

type RecId struct {
  Id string `json:"id"`
}
