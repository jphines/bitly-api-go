package bitly_api

import (
        "net/url"
        "fmt"
        "net/http"
        "io/ioutil"
        "github.com/bitly/go-simplejson"
)

type Connection struct {
  host string
  sslHost string
  login string
  apiKey string
  accessToken string
  secret string
  userAgent string
}

func NewConnection (login string, apiKey string, accessToken string, secret string) *Connection{
  c := new(Connection)
  c.host = "api.bit.ly"
  c.sslHost = "api-ssl.bit.ly"
  c.login = login
  c.apiKey = apiKey
  c.accessToken = accessToken
  c.secret = secret
  c.userAgent = "Go/bitly-api"
  return c
}

func (c *Connection) Shorten (uri string, x_login string, x_apiKey string, preferredDomain string) (map[string]interface{}, error) {
    params := url.Values{}
    params.Set("uri", uri)
    if preferredDomain != "" {
      params.Set("domain", preferredDomain)
    }
    if x_login != "" {
      params.Set("x_login", x_login)
      params.Set("x_apiKey", x_apiKey)
    }
      data, err := c.call(c.host, "v3/shorten", params)
      if err != nil {
        return nil, err
      }
      return data["data"].(map[string]interface{}), nil
}

func (c *Connection) call (host string, endpoint string, params url.Values) (map[string]interface{}, error) {

    var scheme string
    
    if c.accessToken != "" {
        scheme = "https"
        params.Set("access_token", c.accessToken)
        host = "api-ssl.bit.ly"
    } else {
        scheme = "http"
        params.Set("login", c.login)
        params.Set("apiKey", c.apiKey)
    }

    // if c.secret != "" {
        // params.Set("signature", generateSignature(params, c.secret))
    // }
    request_url := fmt.Sprintf("%s://%s/%s?%s", scheme, host, endpoint, params.Encode())

    
    http_client := &http.Client{}
    request, err := http.NewRequest("GET", request_url, nil)
    if err != nil {
        return nil, err
    }

    request.Header.Add("User-agent", c.userAgent)
    response, err := http_client.Do(request)
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }
    js, err := simplejson.NewJson([]byte(contents))
    if err != nil {
        return nil, err
    }
    return js.Map()
}
