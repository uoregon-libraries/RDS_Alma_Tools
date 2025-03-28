package file

import(
  "testing"
  "fmt"
)

func TestFilename(t *testing.T){
  f := Filename()
  fmt.Println(f)
}
