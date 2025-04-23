package withdraw

import(
  "testing"
  "os"
)

func Test_is_serial(t *testing.T){
  data, err := os.ReadFile(os.Getenv("HOME_DIR") + "/fixtures/response_1743689261042.json")
  if err != nil { t.Errorf("problem reading file") }
  is_serial := Is_serial(data)
  if is_serial != true { t.Errorf("incorrect result") }
}
