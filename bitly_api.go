package bitly_api

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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

type UserLink struct {
	title    string
	note     string
	private  bool
	user_ts  string
	archived string
}

type UserLinkHistory struct {
	createdBefore int
	createdAfter  int
	archived      string
	limit         int
	offset        int
}

type Metrics struct {
	unit     string
	units    int
	tzOffset *int
	rollup   *bool
	limit    int
}

func NewConnection(accessToken string, secret string) *Connection {
	c := new(Connection)
	c.accessToken = accessToken
	c.secret = secret
	return c
}

func NewConnectionOauth(accessToken string, apiKey string, login string, secret string) *Connection {
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

func contains(stringSlice []string, elem string) bool {
	for i := 0; i < len(stringSlice); i++ {
		if stringSlice[i] == elem {
			return true
		}
	}
	return false
}

func constructBasicParams(arg string) url.Values {
	params := url.Values{}
	atype := hashOrUrl(arg)
	params.Set(atype, arg)
	return params
}

func convertValueToString(field reflect.Value) (string, bool) {
	switch field.Kind() {
	case reflect.Float64:
		if field.Float() != 0.0 {
			return fmt.Sprintf("%f", field.Float()), true
		}
	case reflect.Bool:
		return fmt.Sprintf("%t", field.Bool()), true
	case reflect.Int:
		if field.Int() != 0 {
			return fmt.Sprintf("%d", field.Int()), true
		}
	case reflect.String:
		if field.String() != "" {
			return fmt.Sprintf("%s", field.String()), true
		}
	case reflect.Ptr:
		if field.Pointer() != 0 {
			field := field.Elem()
			if field.Kind() == reflect.Int {
				return fmt.Sprintf("%d", field.Int()), true
			} else if field.Kind() == reflect.Bool {
				return fmt.Sprintf("%t", field.Int()), true
			}
		}
	}
	return "", false
}

func constructParams(typeStruct interface{}) (url.Values, error) {
	params := url.Values{}
	avalue := reflect.ValueOf(typeStruct)
	atype := reflect.TypeOf(typeStruct)

	if atype.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("Invalid Type %s", atype.Kind()))
	}

	for i := 0; i < atype.NumField(); i++ {
		fieldType := atype.Field(i)
		fieldValue := avalue.Field(i)
		if value, OK := convertValueToString(fieldValue); OK {
			params.Set(fieldType.Name, value)
		}
	}
	return params, nil
}

func constructMetricParams(metrics Metrics) (url.Values, error) {
	allowed_units := []string{"minute", "hour", "day", "week", "mweek", "month"}
	if metrics.unit != "" && metrics.units != 0 {
		if !contains(allowed_units, metrics.unit) {
			return nil, errors.New("Invalid unit")
		}
	}
	if metrics.tzOffset != nil {
		if *metrics.tzOffset != 0 && (*metrics.tzOffset <= -12 || *metrics.tzOffset >= 12) {
			return nil, errors.New("Invalid tzOffset")
		}
	}
	params, err := constructParams(metrics)
	if err != nil {
		return nil, err
	}
	return params, nil
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

func (c *Connection) Info(arg string) (map[string]interface{}, error) {
	return c.call("info", constructBasicParams(arg), true)
}

func (c *Connection) LinkEncodersCount(link string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("link", link)
	return c.call("link/encoders_count", params, false)
}

func (c *Connection) LinkClicks(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/clicks", params, false)
}

func (c *Connection) LinkReferrersByDomain(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/referrers_by_domain", params, false)
}

func (c *Connection) LinkReferrers(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/referrers", params, false)
}

func (c *Connection) LinkShares(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/shares", params, false)
}

func (c *Connection) LinkCountries(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/countries", params, false)
}

func (c *Connection) LinkInfo(link string, metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	params.Set("link", link)
	return c.callOauth2("link/info", params, false)
}

func (c *Connection) LinkContent(link string, contentType string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("link", link)
	params.Set("content_type", contentType)

	return c.callOauth2("link/content", params, false)
}

func (c *Connection) LinkCategory(link string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("link", link)
	return c.callOauth2("link/category", params, false)
}

func (c *Connection) LinkLocation(link string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("link", link)
	return c.callOauth2("link/location", params, false)
}

func (c *Connection) LinkSocial(link string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("link", link)
	return c.callOauth2("link/social", params, false)
}

func (c *Connection) UserClicks(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/clicks", params, false)
}

func (c *Connection) UserCountries(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/countries", params, false)
}

func (c *Connection) UserPopularLinks(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/popular_links", params, false)
}

func (c *Connection) UserReferrers(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/referrers", params, false)
}

func (c *Connection) UserReferringDomains(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/referring_domains", params, false)
}

func (c *Connection) UserShareCounts(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/share_counts", params, false)
}

func (c *Connection) UserShareCountsByType(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/share_counts_by_share_type", params, false)
}

func (c *Connection) UserShortenCounts(metrics Metrics) (map[string]interface{}, error) {
	params, err := constructMetricParams(metrics)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/shorten_counts", params, false)
}

func (c *Connection) UserLinkLookup(uri string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("url", uri)
	return c.call("user/link_lookup", params, true)
}

func (c *Connection) UserLinkEdit(link string, edit string, userLink UserLink) (map[string]interface{}, error) {
	if link == "" || edit == "" {
		return nil, errors.New("UserLinkEdit Missing Args")
	}

	params, err := constructParams(userLink)
	if err != nil {
		return nil, err
	}

	params.Set("link", link)
	params.Set("edit", edit)
	return c.callOauth2("user/link_edit", params, true)
}

func (c *Connection) UserLinkSave(longUrl string, userLink UserLink) (map[string]interface{}, error) {
	if longUrl == "" {
		return nil, errors.New("UserLinkSave Missing Args")
	}
	if userLink.archived != "" {
		return nil, errors.New("UserLinkSave does not support archive")
	}

	params, err := constructParams(userLink)
	if err != nil {
		return nil, err
	}
	params.Set("longUrl", longUrl)

	return c.callOauth2("user/link_save", params, true)
}

func (c *Connection) UserLinkHistory(userHistory UserLinkHistory) (map[string]interface{}, error) {
	params, err := constructParams(userHistory)
	if err != nil {
		return nil, err
	}
	return c.callOauth2("user/link_history", params, true)
}

func (c *Connection) callOauth2(endpoint string, params url.Values, array_wrapped bool) (map[string]interface{}, error) {
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

	api_error := ""
	var data map[string]interface{}
	if array_wrapped {
		if strings.Contains(endpoint, "/") {
			split := strings.Split(endpoint, "/")
			endpoint = split[len(split)-1]
		}
		if endpoint == "link_save" || endpoint == "link_edit" {
			data, err = js.Get("data").Get(endpoint).Map()
			api_error, _ = js.Get("status_txt").String()
		} else {
			data, err = js.Get("data").Get(endpoint).GetIndex(0).Map()
			api_error, _ = js.Get("status_txt").String()
		}
	} else {
		data, err = js.Get("data").Map()
		api_error, _ = js.Get("status_txt").String()
	}
	if api_error != "OK" {
		return nil, errors.New(api_error)
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}
