package bitly_api

import (
        "net/url"
        "fmt"
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

func (c *Connection) shorten(params url.Values) (map[string]interface{}, error){
    return c.call("shorten", params, false)
}

func (c *Connection) ExpandHash(hash string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("hash", hash)
  return c.expand(params)
}

func (c *Connection) ExpandShortUrl(shortUrl string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("shortUrl", shortUrl)
  return c.expand(params)

}

func (c *Connection) expand (params url.Values) (map[string]interface{}, error) {
  return c.call("expand", params, true)
}

func (c *Connection) ClicksHash (hash string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("hash", hash)
  return c.clicks(params)
}

func (c *Connection) ClicksShortUrl(shortUrl string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("shortUrl", shortUrl)
  return c.clicks(params)
}

func (c *Connection) clicks (params url.Values) (map[string]interface{}, error) {
  return c.call("clicks", params, true)
}

func (c *Connection) ClicksByDayHash(hash string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("hash", hash)
  return c.clicksByDay(params)
}

func (c *Connection) ClicksByDayShortUrl(shortUrl string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("shortUrl", shortUrl)
  return c.clicksByDay(params)
}

func (c *Connection) clicksByDay (params url.Values) (map[string]interface{}, error){
  return c.call("clicks_by_day", params, true)
}

func (c *Connection) ClicksByMinuteHash(hash string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("hash", hash)
  return c.clicksByMinute(params)
}

func (c *Connection) ClicksByMinuteShortUrl(shortUrl string) (map[string]interface{}, error) {
  params := url.Values{}
  params.Set("shortUrl", shortUrl)
  return c.clicksByMinute(params)
}

func (c *Connection) clicksByMinute (params url.Values) (map[string]interface{}, error){
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
