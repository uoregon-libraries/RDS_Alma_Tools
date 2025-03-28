package utils

import(
  "net/http"
  "time"
  "github.com/labstack/echo/v4"
  "log"
  "errors"
)

//expires value is in minutes
func WriteCookie(c echo.Context, expires int, name string, value string) {
  cookie := new(http.Cookie)
  cookie.Name = name
  cookie.Value = value
  cookie.Path = "/"
  cookie.Expires = time.Now().Add(time.Duration(expires) * time.Minute)
  c.SetCookie(cookie)
}

func FetchCookieVal(c echo.Context, name string) (string, error) {
  cookie, err := c.Cookie(name)
  if err != nil { log.Println(err); return "", errors.New("unable to fetch cookie") }

  return cookie.Value, nil
}


