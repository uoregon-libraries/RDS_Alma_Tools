package withdraw

import (
  "testing"
  "fmt"
  "net/url"
  "os"
  "errors"
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

func TestFinalReport(t *testing.T){
  e := Eligible{ SerialRequiresAction: true, Locations: []string{"kshort"} }
  eligibleLists := map[string]Eligible{}
  eligibleLists["123123123123"] = e
  Final_Report("test-final", eligibleLists)
  path := os.Getenv("REPORT_DIR") + "/" + "test-final"
  _, err := os.Stat(path)
  if err != nil {
    if errors.Is(err, os.ErrNotExist) { t.Errorf("did not write report") }
  }
}
