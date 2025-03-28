package withdraw

import (
  "testing"
  "fmt"
  "net/url"
)

func TestBuildItemLink(t *testing.T){
  mmsId := "12345"
  holdingId := "67891"
  pid := "5432"
  link := BuildItemLink(mmsId, holdingId, pid)
  correct := fmt.Sprintf("/almaws/v1/bibs/%s/holdings/%s/items/%s", mmsId, holdingId, pid)
  parsed,_ := url.Parse(link)
  if parsed.Path != correct { t.Errorf("incorrect path") }
}
