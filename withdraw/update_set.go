package withdraw

import (
  "encoding/json"
  "rds_alma_tools/connect"
  "os"
  "log"
)

//setcontent either BIB_MMS or ITEM
func UpdateSet(filename string, setname string, setcontent string, eligibleList []string)error{
  setid := os.Getenv(setname)
  params := []string{ "op=replace_members", ApiKey() }
  url := BaseUrl() + "conf/sets/" + setid
  set := InitSet(setname, setcontent)
  set = SetMembers(set, eligibleList)
  body, err := json.Marshal(set)
  if err != nil {}
  _, err = connect.Put(url, params, string(body))
  if err != nil { log.Println(err); return err }
  return nil 
}

func SetMembers(set Set, eligibleList []string)Set{
  for _,v := range eligibleList {
    set.Members.Member = append(set.Members.Member, BibId{ Id: v })
  }
  return set
}

func InitSet(name string, content string) Set{
  var set = Set{Name: name, Type: Val{Value: "LOGICAL"}, Content: Val{Value: content}, Query: Val{Value: ""}, Members: MemberArr{ Member: []BibId{} }}
  return set
}

type Set struct {
  Name string
  Type Val
  Content Val
  Query Val
  Members MemberArr
}

type Val struct {
  Value string
}

type MemberArr struct {
  Member []BibId
}

type BibId struct {
  Id string
}
