package watcher

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/site"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

var (
    server *httptest.Server
)

func init() {
    core.NewApp(":3000", "sqlite3", "/tmp/test.db", true)
    site.Setup("../site/templates")

    site.Migrate()

    server = httptest.NewServer(core.App.Router)

    // Replace Alerter with mock
    Shout = &MockAlerter{}
    clear()
}

func clear() {
    core.App.DB.Delete(&core.Subject{})
    core.App.DB.Delete(&core.Audit{})
}

func TestWatchWithSuccessSubject(t *testing.T) {
    clear()
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/health",
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
    clear()
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/health",
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
    clear()
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/test/fakse",
        ResponseLimit: 5,
        Name:          "TestCrit",
    }

    core.App.DB.Create(s)

    Watch()

    var o core.Subject
    core.App.DB.Where("ext_id = ?", s.ExtId).Find(&o)

    assert.Equal(t, CRITICAL, o.Status)
}

func TestAdvancedHealthEndpoint(t *testing.T) {
    clear()
    s := &core.Subject{
        Domain:        server.URL,
        PingURI:       "/health",
        ResponseLimit: 5,
        Name:          "TestHealth",
        Advanced:      false,
    }

    core.App.DB.Create(s)

    Watch()

    var o core.Subject

    core.App.DB.Where("ext_id = ?", s.ExtId).Preload("Audits").Find(&o)

    assert.Equal(t, OK, o.Status)
    assert.NotEqual(t, 0, o.Audits[0].Memory)
    assert.NotEqual(t, 0, o.Audits[0].CPU)
}

func TestFailToConnectDuringTest(t *testing.T) {
    clear()
    s := &core.Subject{
        Domain:        "http://rand.rand",
        PingURI:       "/health",
        ResponseLimit: 5,
        Name:          "TestFailure",
        Advanced:      false,
    }

    core.App.DB.Create(s)

    // Update the HTTPClient, 10 seconds is nuts for a test
    HttpClient = http.Client{
        Timeout: time.Duration(1 * time.Second),
    }

    Watch()

    var o core.Subject

    core.App.DB.Where("ext_id = ?", s.ExtId).Preload("Audits").Find(&o)

    assert.Equal(t, CRITICAL, o.Status)
    assert.Equal(t, CRITICAL, o.Audits[0].Status)
}
