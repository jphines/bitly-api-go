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
    contents, err := c.call(c.host, "v3/shorten", params)
     if err != nil {
      return nil, err
    }

    js, err := simplejson.NewJson([]byte(contents))
    if err != nil {
      return nil, err
    }

    data, err := js.Map()
    if err != nil {
      return nil, err
    }

    return data["data"].(map[string]interface{}), nil
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

func (c *Connection) expand(params url.Values) (map[string]interface{}, error){
  contents, err := c.call(c.host, "v3/expand", params)
  if err != nil {
      return nil, err
  }

  js, err := simplejson.NewJson([]byte(contents))
  if err != nil {
      return nil, err
  }
  data, err := js.Get("data").Get("expand").GetIndex(0).Map()
  if err != nil {
    return nil, err
  }
  return data, nil
}

func (c *Connection) call (host string, endpoint string, params url.Values) ([]byte, error) {

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
    
    return ioutil.ReadAll(response.Body)
  }
