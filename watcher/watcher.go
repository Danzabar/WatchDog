package watcher

import (
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "net/http"
    "time"
)

var (
    HttpClient http.Client
    Shout      Alerter
)

const (
    OK         = "ok"
    CRITICAL   = "critical"
    DEGREDATED = "degredated"
)

func init() {
    HttpClient = http.Client{
        Timeout: time.Duration(10 * time.Second),
    }

    Shout = NewPushBullet()
}

func Watch() {
    var s []core.Subject

    core.App.Log.Debug("Starting watcher...")
    core.App.DB.Find(&s)

    for _, v := range s {
        core.App.Log.Debugf("Checking %s", v.Domain)

        a := CheckStatus(v, &core.Audit{SubjectId: v.Model.ID})
        a.Status = AnalyseStatus(a.Result, a.ResponseTime, v)
        v.Status = a.Status

        core.App.DB.Save(a)
        core.App.DB.Save(&v)
        core.App.Log.Debugf("Checked %s", v.Domain)
    }
}

func AnalyseStatus(r bool, t float64, s core.Subject) string {
    if !r {
        // If we go to critical, we want an alert for this
        Shout.SendAlert(
            fmt.Sprintf("%s domain has entered critical status", s.Domain),
            fmt.Sprintf("CRITICAL: %s", s.Domain),
        )
        return CRITICAL
    }

    // Why 2? I don't really know, but 2 seconds seems
    // like a long time for a ping endpoint
    if t > 2 {
        // Degredation? Yes please
        Shout.SendAlert(
            fmt.Sprintf("%s domain has entered degredated status", s.Domain),
            fmt.Sprintf("DEGREDATION: %s", s.Domain),
        )
        return DEGREDATED
    }

    // If we were previously at a different status, and we are ok
    // We want to know
    if s.Status != OK {
        Shout.SendAlert(
            fmt.Sprintf("%s domain is now running OK", s.Domain),
            fmt.Sprintf("OK: %s", s.Domain),
        )
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
