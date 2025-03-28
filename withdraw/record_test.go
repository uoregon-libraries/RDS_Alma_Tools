package withdraw

import(
  "testing"
)

func TestLineMap(t *testing.T){
  line := "9984898401852	XBox 360	12345678	22274069860001852	23193212440001852	35025040997286	Item not in place	Science	sgames	fake public note	toggled missing status from technical migration. was breaking bookings - SDG	STATUS2: r|ICODE2: p|I TYPE2: 77|LOCATION: orvng|RECORD #(ITEM)2: i45612675	NOTE(ITEM): serial number: 118381693005	Status: r - IN REPAIR, 2018/1/26 toggled missing status from technical migration. was breaking bookings - SDG	fake_retention_note"
  lineMap := LineMap(line)
  if lineMap["retention_note"] != "fake_retention_note" { t.Errorf("new library value is wrong") }
}

func TestOclcSelect(t *testing.T){
  vals := []string{"(OCoLC)ocm12345678", "(OrU)b36609079-01alliance_uo"}
  oclc := OclcSelect(vals)
  if oclc != "12345678" { t.Errorf("wrong oclc selected") }
}

func TestLineMap2(t *testing.T){
  line := "99901146446201852	UO test 1	4270644	22474038060001852	23474038030001852	ALMA615688	Item in place	Knight	kgen						"
  lineMap := LineMap(line)
  if lineMap["library"] != "Knight" { t.Errorf("new library value is wrong") }
}
