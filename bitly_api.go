package bitly_api

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
  "errors"
)

type Json struct {
	data interface{}
}

type Connection struct {
	login       string
	apiKey      string
	accessToken string
	secret      string
}

func NewConnection(accessToken string, secret string) *Connection {
	c := new(Connection)
	c.accessToken = accessToken
	c.secret = secret
	return c
}

func NewConnectionOauth(accessToken, string, apiKey string, login string, secret string) *Connection {
	c := NewConnection(accessToken, secret)
	c.apiKey = apiKey
	c.login = login
	return c
}

func (c *Connection) ShortenWithDomain(uri string, preferredDomain string) (map[string]interface{}, error) {
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

func hashOrUrl(mystery string) string {
	var solved string
	if strings.Contains(mystery, "/") {
		solved = "shortUrl"
	} else {
		solved = "hash"
	}
	return solved
}

func constructBasicParams(arg string) (url.Values) {
	params := url.Values{}
	atype := hashOrUrl(arg)
	params.Set(atype, arg)
  return params
}

func constructLinkParams(arg string) (url.Values) {
	params := url.Values{}
  params.Set("link", arg)
  return params
}

func (c *Connection) shorten(params url.Values) (map[string]interface{}, error) {
	return c.call("shorten", params, false)
}

func (c *Connection) Expand(arg string) (map[string]interface{}, error) {
	return c.call("expand", constructBasicParams(arg), true)
}

func (c *Connection) Clicks(arg string) (map[string]interface{}, error) {
	return c.call("clicks", constructBasicParams(arg), true)
}

func (c *Connection) ClicksByDay(arg string) (map[string]interface{}, error) {
	return c.call("clicks_by_day", constructBasicParams(arg), true)
}

func (c *Connection) ClicksByMinute(arg string) (map[string]interface{}, error) {
	return c.call("clicks_by_minute", constructBasicParams(arg), true)
}

func (c *Connection) Referrers(arg string) (map[string]interface{}, error) {
  return c.call("referrers", constructBasicParams(arg), true)
}

func (c *Connection) Info (arg string) (map[string]interface{}, error) {
	return c.call("info", constructBasicParams(arg), true)
}

func (c *Connection) LinkEncodersCount(arg string) (map[string]interface{}, error) {
  return c.call("link/encoders_count", constructLinkParams(arg), false)
}

func (c *Connection) callOath2(endpoint string, params url.Values, array_wrapped bool) (map[string]interface{}, error) {
    if c.accessToken == "" {
        return nil, errors.New(fmt.Sprintf("This endpoint %s requires Oauth", endpoint))
    }
    return c.call(endpoint, params, array_wrapped)
}

func (c *Connection) call(endpoint string, params url.Values, array_wrapped bool) (map[string]interface{}, error) {

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
  fmt.Println(request_url)

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
