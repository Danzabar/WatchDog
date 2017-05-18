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

func clear() {
    core.App.DB.Delete(&core.Subject{})
    core.App.DB.Delete(&core.Audit{})
}

func TestPostSubjectSuccess(t *testing.T) {
    s := &core.Subject{Name: "Test", Domain: server.URL, PingURI: "/ping"}
    js, _ := json.Marshal(s)

    clear()

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

func TestGetSubject(t *testing.T) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/subject", server.URL), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSingleSubjectSuccess(t *testing.T) {
    clear()
    s := &core.Subject{Name: "TestGetSingle"}
    core.App.DB.Save(s)

    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/subject/%s", server.URL, s.ExtId), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSingeSubjectNotFound(t *testing.T) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/subject/test", server.URL), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteSubjectSuccess(t *testing.T) {
    clear()
    s := &core.Subject{Name: "Testing"}
    core.App.DB.Save(s)

    req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/subject/%s", server.URL, s.ExtId), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteSubjectNotFound(t *testing.T) {
    req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/subject/fake", server.URL), nil)
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
