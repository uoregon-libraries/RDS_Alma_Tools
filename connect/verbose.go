package connect

import(
  "net/http"
  "net/http/httputil"
  "log"
)

func RequestDump(verbose string, req *http.Request){
  if verbose == "true" {
    reqdump, err := httputil.DumpRequest(req, true)
    if err != nil { log.Println(err) } else {
      log.Printf("REQUEST:\n%s", string(reqdump)) }
  }
}

func ResponseDump(verbose string, response *http.Response){
  if verbose == "true" {
    respdump, err := httputil.DumpResponse(response, true)
    if err != nil { log.Println(err) } else {
      log.Printf("RESPONSE:\n%s", string(respdump)) }
  }
}
