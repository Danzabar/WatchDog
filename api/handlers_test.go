package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/site"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

var server *httptest.Server

func init() {
    core.NewApp(":3000", "sqlite3", "/tmp/test.db", false)
    Setup()

    site.Migrate()
    server = httptest.NewServer(core.App.Router)
}

func TestPostSubjectSuccess(t *testing.T) {
    s := &core.Subject{Name: "Test"}
    js, _ := json.Marshal(s)

    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/subject", server.URL), bytes.NewReader(js))
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostSubjectBadJson(t *testing.T) {
    js := []byte(`{"test":}`)

    req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/subject", server.URL), bytes.NewReader(js))
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
