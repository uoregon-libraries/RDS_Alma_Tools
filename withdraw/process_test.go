package withdraw

import (
  "testing"
  "slices"
  "fmt"
  "net/url"
)

func TestExtractEligibles( t *testing.T){
  list := map[string][]bool{}
  list["a"] = []bool{true, true}
  list["b"] = []bool{true, false}
  list["c"] = []bool{false, true}
  newlist := ExtractEligibles(list, 0)
  if !slices.Equal(newlist, []string{"a", "b"}) { t.Errorf("incorrect extraction") }
}

func TestBuildItemLink(t *testing.T){
  mmsId := "12345"
  holdingId := "67891"
  pid := "5432"
  link := BuildItemLink(mmsId, holdingId, pid)
  correct := fmt.Sprintf("/almaws/v1/bibs/%s/holdings/%s/items/%s", mmsId, holdingId, pid)
  parsed,_ := url.Parse(link)
  if parsed.Path != correct { t.Errorf("incorrect path") }
}
