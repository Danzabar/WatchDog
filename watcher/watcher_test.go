package watcher

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/site"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

var (
    server *httptest.Server
)

// Test ping routes
func RouteForSuccess(w http.ResponseWriter, r *http.Request) {
    core.WriteResponseHeader(w, 200)
}

func init() {
    core.NewApp(":3000", "sqlite3", "/tmp/test.db", true)
    site.Setup("../site/templates")

    site.Migrate()

    // Add test routes
    core.App.Router.HandleFunc("/test/success", RouteForSuccess)

    server = httptest.NewServer(core.App.Router)

    // Replace Alerter with mock
    Shout = &MockAlerter{}

    core.App.DB.Delete(&core.Subject{})
}

func TestWatchWithSuccessSubject(t *testing.T) {
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/test/success",
        ResponseLimit: 5,
        Name:          "TestSuccess",
    }

    core.App.DB.Create(s)

    Watch()

    var o core.Subject
    core.App.DB.Where("ext_id = ?", s.ExtId).Find(&o)

    assert.Equal(t, OK, o.Status)
}

func TestWatchWithDegredation(t *testing.T) {
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/test/success",
        ResponseLimit: 0.00001,
        Name:          "TestDeg",
    }

    core.App.DB.Create(s)

    Watch()

    var o core.Subject
    core.App.DB.Where("ext_id = ?", s.ExtId).Find(&o)

    assert.Equal(t, DEGREDATED, o.Status)
}

func TestWatchWithCritical(t *testing.T) {
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/test/fakse",
        ResponseLimit: 5,
        Name:          "TestDeg",
    }

    core.App.DB.Create(s)

    Watch()

    var o core.Subject
    core.App.DB.Where("ext_id = ?", s.ExtId).Find(&o)

    assert.Equal(t, CRITICAL, o.Status)
}
