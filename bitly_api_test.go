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
	_, err := bitly.Shorten(testUrl)
	if err != nil {
		t.Fatalf("bitly Shorten returned an err %s", err)
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

  _, err := bitly.Info(testUrl)
	if err != nil {
		t.Fatalf("bitly Info returned an error %s", err)
	}
}

func TestLinkEncodersCount(t *testing.T) {
	bitly := getConnection(t)

  _, err := bitly.LinkEncodersCount(testUrl)
	if err != nil {
		t.Fatalf("bitly links_encoders_count returned an error %s", err)
	}
}

func TestUserLink(t *testing.T) {
  bitly := getConnection(t)

  _, err := bitly.UserLinkLookup(testUrl)
	if err != nil {
	  t.Fatalf("bitly UserLinkLookup returned an error %s", err)
	}
  
  _, err = bitly.UserLinkSave(longUrl, UserLink{private:true})
  if err != nil && err.Error() != "LINK_ALREADY_EXISTS" {
    t.Fatalf("bitly user/link_save returned an expected result %s", err)
  }

  _, err = bitly.UserLinkEdit(linkUrl, "title", UserLink{title:"New Title"})
  if err != nil {
    t.Fatalf("bitly user/link_edit returned an expected result %s", err)
  }

  _, err = bitly.UserLinkHistory(UserLinkHistory{archived:"on"})
  if err != nil {
    t.Fatalf("bitly user/link_history returned an expected result %s", err)
  }
}

func TestLinkMetrics(t *testing.T) {
  bitly := getConnection(t)
  _, err := bitly.LinkClicks(testUrl, Metrics{limit:1})
   if err != nil {
    t.Fatalf("bitly link/clicks returned an error %s", err)
  }
}
