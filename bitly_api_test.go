package bitly_api

import (
	"os"
	"testing"
)

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
	if bitly == nil {
		t.Fatalf("bitly connection returned nil")
	}
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
	if bitly == nil {
		t.Fatalf("bitly connection returned nil")
	}
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
	if bitly == nil {
		t.Fatalf("bitly connection returned nil")
	}
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
		t.Fatalf("bitly ClicksHash did not return NOT_FOUND", err)
	}

}
