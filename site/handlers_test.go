package site

import (
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

var (
    server *httptest.Server
)

func init() {
    core.NewApp(":3000", "sqlite3", "/tmp/test.db")
    Setup()

    server = httptest.NewServer(core.App.Router)
}

// Test that the Ping endpoint works
func TestPingIsOk(t *testing.T) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/ping", server.URL), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Test that we get a 200 from stats
func TestStatsIsOk(t *testing.T) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/status", server.URL), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
