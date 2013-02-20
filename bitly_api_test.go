package bitly_api

import (
        "os"
        "testing"
)


func getConnection(t *testing.T) *Connection {

    token := "BITLY_API_TOKEN"

    BITLY_API_TOKEN := os.Getenv(token)
    if BITLY_API_TOKEN == "" {
        t.Fatalf(token + " not found")
        return nil
    }
    accessToken := BITLY_API_TOKEN
    return &Connection{accessToken:accessToken}
}

func TestApi(t *testing.T) {
    bitly := getConnection(t)
    if bitly == nil {
        t.Fatalf("bitly connection returned nil")
    }
    testUrl := "http://google.com/"
    data, err := bitly.Shorten(testUrl, "", "", "")
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
