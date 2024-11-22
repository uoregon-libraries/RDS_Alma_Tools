package withdraw

import(
  "github.com/labstack/echo/v4"
  "rds_alma_tools/connect"
  "log"
  "net/http"
  "github.com/tidwall/sjson"
  "os"
  "bufio"
  //"github.com/tidwall/gjson"
  "errors"
  "strings"
  "net/url"
  "strconv"
  "time"
  "fmt"
)

func ProcessHandler(c echo.Context)(error){
  //get uploaded file
  err := UpdateItems(c)
  if err != nil { log.Println(err); return c.String(http.StatusBadRequest, "Could not process items") }
  
  //next steps...

  return c.String(http.StatusOK, "")
}

func BuildItemLink(mmsId string, holdingId string, pid string)string{
  _url,_ := url.Parse(BaseUrl())
  _url = _url.JoinPath("bibs", mmsId, "holdings", holdingId, "items", pid)
  return _url.String()
}

func UpdateItems(c echo.Context)error{
  file, _ := c.FormFile("file")
  src, err := file.Open()
  if err != nil { log.Println(err); return errors.New("Unable to open file") }
  defer src.Close()

  //process the data for updating item records
  scanner := bufio.NewScanner(src)
  for scanner.Scan(){
    line := scanner.Text()
    itemRec, err := UpdateItem(c.Param("loc_type"), line)
    if err != nil {}
    fmt.Println(string(itemRec))
    //err := connect.Post(url, params, itemRec)
  }
  return nil
}

func TimeNow()time.Time{
  loc, _ := time.LoadLocation("America/Los_Angeles")
  t := time.Now().In(loc)
  return t
}
func FiscalYear(t time.Time)string{
  m := t.Format("01")
  y := t.Format("2006")
  mInt, _ := strconv.Atoi(m)
  if mInt > 6 { return y } else {
    yInt, _ := strconv.Atoi(y)
    return strconv.Itoa(yInt-1)
  }
}

func UpdateItem(newLocType string, data string)([]byte, error){
  arr := strings.Split(data, "\t")
  url := BuildItemLink(arr[0], arr[2], arr[3])
  params := []string{ ApiKey() }
  itemRec,_ := connect.Get(url, params)
  libMap := LoadMap()
  newLibVal := libMap[Key{arr[6],newLocType,"value"}]
  newLibDesc := libMap[Key{arr[6],newLocType,"desc"}]
  internalNote3 := arr[12] + "|WD FY" + FiscalYear(TimeNow())
  newStatusDesc := "missing"
  newStatusVal := ""
  //using sjson insert new library, status, append note
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.value", newLibVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.library.desc", newLibDesc)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.internal_note_3", internalNote3)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.base_status.value", newStatusVal)
  itemRec,_ = sjson.SetBytes(itemRec, "item_data.base_status.desc", newStatusDesc)
  return itemRec, nil
}

type Key struct{
  Libname, Libtype, Proptype string
}

func LoadMap()map[Key]string{
  libmap := map[Key]string{}
  homedir := os.Getenv("HOME_DIR")
  data,_ := os.Open(homedir + "/withdraw/library_map.txt")
  fileScanner := bufio.NewScanner(data)
  fileScanner.Split(bufio.ScanLines)
  for fileScanner.Scan(){
    arr := strings.Split(fileScanner.Text(), "\t")
    libmap[Key{arr[0], arr[1], arr[2]}] = arr[3]
  }
  return libmap
}

func ApiKey()string{
  key := os.Getenv("ALMA_KEY")
  return "apikey=" + key
}
func BaseUrl()string{
  return os.Getenv("ALMA_URL")
}
