package watcher

import (
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "net/http"
    "time"
)

var HttpClient http.Client

const (
    OK         = "ok"
    CRITICAL   = "critical"
    DEGREDATED = "degredated"
)

func init() {
    HttpClient = http.Client{
        Timeout: time.Duration(10 * time.Second),
    }
}

func Watch() {
    var s []core.Subject

    core.App.Log.Debug("Starting watcher...")
    core.App.DB.Find(&s)

    for _, v := range s {
        core.App.Log.Debugf("Checking %s", v.Domain)

        a := CheckStatus(v, &core.Audit{SubjectId: v.Model.ID})
        a.Status = AnalyseStatus(a.Result, a.ResponseTime)
        v.Status = a.Status

        core.App.DB.Save(a)
        core.App.DB.Save(&v)
        core.App.Log.Debugf("Checked %s", v.Domain)
    }
}

func AnalyseStatus(s bool, t float64) string {
    if !s {
        return CRITICAL
    }

    if t > 4 {
        return DEGREDATED
    }

    return OK
}

// Checks the Status of a service/website using a HTTP "ping"
func CheckStatus(s core.Subject, a *core.Audit) *core.Audit {
    a.Result = false
    ts := time.Now()
    req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", s.Domain, s.PingURI), nil)

    if err != nil {
        core.App.Log.Error(err)
        return a
    }

    resp, err := HttpClient.Do(req)

    if err != nil {
        core.App.Log.Error(err)
        a.ResponseTime = time.Since(ts).Seconds()
        return a
    }

    if resp.StatusCode == http.StatusOK {
        a.Result = true
        a.ResponseTime = time.Since(ts).Seconds()
    }

    return a
}
