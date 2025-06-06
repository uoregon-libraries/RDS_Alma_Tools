package withdraw

import(
  "testing"
  "os"
)

func Test_IsSerial(t *testing.T){
  data, err := os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/response_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  is_serial := IsSerial(data)
  if is_serial != true { t.Errorf("incorrect result") }
}

func Test_IsBoundWith(t *testing.T){
  //case 0
  data, err := os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/response_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  boundwith, biblist := IsBoundWith(data)
  if boundwith != false { t.Errorf("incorrect boundwith case 0") }
  if biblist != "" { t.Errorf("incorrect boundwithmultiple case 0") }

  //case 1
  data, err = os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/fakeBW1_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  boundwith, biblist = IsBoundWith(data)
  if boundwith != true { t.Errorf("incorrect boundwith case 1") }
  if biblist != "123123123123" { t.Errorf("incorrect boundwithmultiple case 1") }

  //case 2
  data, err = os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/fakeBW2_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  boundwith, biblist = IsBoundWith(data)
  if boundwith != true { t.Errorf("incorrect boundwith case 2") }
  if biblist != "" { t.Errorf("incorrect boundwithmultiple case 2") }

  //case 3
  data, err = os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/fakeBW3_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  boundwith, biblist = IsBoundWith(data)
  if boundwith != true { t.Errorf("incorrect boundwith case 3") }
  if biblist != "" { t.Errorf("incorrect boundwithmultiple case 3") }

}

func Test_Subfield(t *testing.T){
  data, err := os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/fakeBW2_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  marc_string := ExtractJsonField("anies.0", data)
  sub_elt := Subfield("962", "a", marc_string)
  if sub_elt != "BoundWithRecord" { t.Errorf("incorrect response for subfield") }
}
