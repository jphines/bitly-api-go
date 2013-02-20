package bitly_api

import (
        "net/url"
        "fmt"
        "strings"
        "net/http"
        "io/ioutil"
        "github.com/bitly/go-simplejson"
)

type Json struct {
  data interface{}
}

type Connection struct {
  login string
  apiKey string
  accessToken string
  secret string
}

func NewConnection (accessToken string, secret string) *Connection{
  c := new(Connection)
  c.accessToken = accessToken
  c.secret = secret
  return c
}

func NewConnectionApiKey(apiKey string, login string, secret string) *Connection {
  c := new(Connection)
  c.apiKey = apiKey
  c.login = login
  c.secret = secret
  return c
}

func (c *Connection) ShortenWithDomain (uri string,  preferredDomain string) (map[string]interface{}, error) {
    params := url.Values{}
    params.Set("uri", uri)
    params.Set("domain", preferredDomain)
    return c.shorten(params)
}

func (c *Connection) Shorten(uri string) (map[string]interface{}, error) {
    params := url.Values{}
    params.Set("uri", uri)
    return c.shorten(params)
}

func hashOrUrl (mystery string) string {
    var solved string
    if strings.Contains(mystery, "/") {
        solved = "shortUrl"
    } else {
        solved = "hash"
    }
    return solved
}

func (c *Connection) shorten(params url.Values) (map[string]interface{}, error){
    return c.call("shorten", params, false)
}

func (c *Connection) Expand(arg string) (map[string]interface{}, error) {
  params := url.Values{}
  atype := hashOrUrl(arg)
  params.Set(atype, arg)
  return c.call("expand", params, true)
}

func (c *Connection) Clicks (arg string) (map[string]interface{}, error) {
  params := url.Values{}
  atype := hashOrUrl(arg)
  params.Set(atype, arg)
  return c.call("clicks", params, true)
}

func (c *Connection) ClicksByDay (arg string) (map[string]interface{}, error) {
  params := url.Values{}
  atype := hashOrUrl(arg)
  params.Set(atype, arg)
  return c.call("clicks_by_day", params, true)
}

func (c *Connection) ClicksByMinute(arg string) (map[string]interface{}, error) {
  params := url.Values{}
  atype := hashOrUrl(arg)
  params.Set(atype, arg)
  return c.call("clicks_by_minute", params, true)
}

func (c *Connection) call (endpoint string, params url.Values, array_wrapped bool) (map[string]interface{}, error) {

    var scheme string
    var host string
    
    if c.accessToken != "" {
        scheme = "https"
        params.Set("access_token", c.accessToken)
        host = "api-ssl.bit.ly"
    } else {
        scheme = "http"
        params.Set("login", c.login)
        params.Set("apiKey", c.apiKey)
        host = "api.bit.ly"
    }

    // if c.secret != "" {
        // params.Set("signature", generateSignature(params, c.secret))
    // }
    request_url := fmt.Sprintf("%s://%s/v3/%s?%s", scheme, host, endpoint, params.Encode())

    
    http_client := &http.Client{}
    request, err := http.NewRequest("GET", request_url, nil)
    if err != nil {
        return nil, err
    }

    request.Header.Add("User-agent", "Go/bitly-api")
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
    
    var data map[string]interface{}
    if array_wrapped {
        data, err = js.Get("data").Get(endpoint).GetIndex(0).Map()
    } else {
        data, err = js.Get("data").Map()
    }
    if err != nil {
        return nil, err
    }

    return data, nil
}
