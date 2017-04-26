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
        core.App.Log.Debugf("Checking %s", v.Name)

        a := CheckStatus(v, &core.Audit{SubjectId: v.Model.ID})
        a.Status = AnalyseStatus(a, v)
        v.Status = a.Status

        core.App.DB.Save(a)
        core.App.DB.Save(&v)
        core.App.Log.Debugf("Checked %s", v.Name)
    }
}

func AnalyseStatus(a *core.Audit, s core.Subject) string {
    if !a.Result {
        // If we go to critical, we want an alert for this
        Shout.SendAlert(
            fmt.Sprintf("%s domain has entered critical status - server responded with a status of %d", s.Domain, a.ResponseStatus),
            fmt.Sprintf("CRITICAL: %s", s.Name),
        )
        return CRITICAL
    }

    if a.ResponseTime > s.ResponseLimit {
        // Degredation? Yes please
        Shout.SendAlert(
            fmt.Sprintf("%s domain has entered degredated status - response time was %d", s.Domain, a.ResponseTime),
            fmt.Sprintf("DEGREDATION: %s", s.Name),
        )
        return DEGREDATED
    }

    // If we were previously at a different status, and we are ok
    // We want to know
    if s.Status != OK {
        Shout.SendAlert(
            fmt.Sprintf("%s domain is now running OK", s.Domain),
            fmt.Sprintf("OK: %s", s.Name),
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

    a.ResponseStatus = resp.StatusCode

    if resp.StatusCode == http.StatusOK {
        a.Result = true
        a.ResponseTime = time.Since(ts).Seconds()
    }

    return a
}
