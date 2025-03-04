package withdraw

import(
  "bufio"
  "math/rand"
  "os"
  "fmt"
  "log"
  "strings"
)

func buildFromFile(filename string)[]string{
  str_arr := []string{}
  f,_ := os.Open(filename)
  defer f.Close()
  scanner := bufio.NewScanner(f)
  for scanner.Scan(){
    arr := strings.Split(scanner.Text(), ",")
    str_arr = append(str_arr, arr...)
  }
  return str_arr
}

func Filename()string{
  homedir := os.Getenv("HOME_DIR")
  nouns := buildFromFile(homedir + "/withdraw/nouns.txt")
  mods := buildFromFile(homedir + "/withdraw/modifiers.txt")
  t := TimeNow()
  y := t.Format("2006")
  m := t.Format("01")
  d := t.Format("02")
  h := t.Format("15")
  mi := t.Format("04")
  s1 := mods[rand.Intn(len(mods))]
  s2 := nouns[rand.Intn(len(nouns))]
  return fmt.Sprintf("%s-%s-%s-%s:%s_%s-%s", y,m,d,h,mi,s1,s2)
}

func WriteReport(filename, message string){
  report_dir := os.Getenv("REPORT_DIR")
  path := report_dir + "/" + filename
  f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  defer f.Close()
  if err != nil { log.Println(err) }
  _, err = fmt.Fprintln(f, message)
  if err != nil { log.Println(err) }
}
