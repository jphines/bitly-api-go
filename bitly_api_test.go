package bitly_api

import (
	"os"
	"testing"
)

var linkUrl string = "http://bitly.com/Xlq5ZH"
var testUrl string = "http://bitly.com/3hQYj"
var longUrl string = "http://google.com"

func getConnection(t *testing.T) *Connection {

	token := "BITLY_ACCESS_TOKEN"

	BITLY_ACCESS_TOKEN := os.Getenv(token)
	if BITLY_ACCESS_TOKEN == "" {
		t.Fatalf(token + " not found")
		return nil
	}
	accessToken := BITLY_ACCESS_TOKEN
  return NewConnection(accessToken, "")
}

func TestApi(t *testing.T) {
	bitly := getConnection(t)
	testUrl := "http://google.com/"
	data, err := bitly.Shorten(testUrl)
	if err != nil {
		t.Fatalf("bitly Shorten returned an err %s", err)
	}

	longUrl := data["long_url"].(string)
	if testUrl != longUrl {
		t.Fatalf("test url != long url from return")
	}

	hash := data["hash"].(string)
	if hash == "" {
		t.Fatalf("hash empty")
	}
}

func TestExpand(t *testing.T) {
	bitly := getConnection(t)
	
  data, err := bitly.Expand("test1_random_fjslfjieljfklsjflkas")
	if err != nil {
		t.Fatalf("bitly Expand returned an error %s", err)
	}
	if data["error"] != "NOT_FOUND" {
		t.Fatalf("bitly Expand did not return NOT_FOUND", err)
	}
}

func TestClicks(t *testing.T) {
	bitly := getConnection(t)
  
  data, err := bitly.Clicks("test1_random_fjslfjieljfklsjflkas")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] != "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}

	data, err = bitly.ClicksByDay("test1_random_fjslfjieljfklsjflkas")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] != "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}

	data, err = bitly.ClicksByMinute("test1_random_fjslfjieljfklsjflkas")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] != "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}
	data, err = bitly.Clicks("http://bit.ly/3hQYj")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] == "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}

	data, err = bitly.ClicksByDay("http://bit.ly/3hQYj")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] == "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}

	data, err = bitly.ClicksByMinute("http://bit.ly/3hQYj")
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["error"] == "NOT_FOUND" {
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND %s", err)
	}
}

func TestInfo(t *testing.T) {
	bitly := getConnection(t)

  data, err := bitly.Info(testUrl)
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["short_url"] != testUrl {
		t.Fatalf("bitly /info returned an unexpected result")
	}
}

func TestLinkEncodersCount(t *testing.T) {
	bitly := getConnection(t)

  data, err := bitly.LinkEncodersCount(testUrl)
	if err != nil {
		t.Fatalf("bitly clicks returned an error %s", err)
	}
	if data["aggregate_link"] != testUrl {
		t.Fatalf("bitly link/encoders_count returned an unexpected result")
	}
}

func TestUserLink(t *testing.T) {
  bitly := getConnection(t)

  data, err := bitly.UserLinkLookup(testUrl)
	if err != nil {
	  t.Fatalf("bitly UserLinkLookup returned an error %s", err)
	}
  if data["url"] != testUrl {
    t.Fatalf("bitly user/link_lookup returned an expected result")
  }

  data, err = bitly.UserLinkSave(longUrl, UserLink{private:true})
  if data["long_url"] != longUrl {
    t.Fatalf("bitly user/link_save returned an expected result %s", data["link"])
  }

  data, err = bitly.UserLinkEdit(linkUrl, "title", UserLink{title:"New Title"})
  if data["link"] != linkUrl {
    t.Fatalf("bitly user/link_edit returned an expected result %s", data["link"])
  }
}

