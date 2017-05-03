package watcher

import (
    "encoding/json"
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/plugins"
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
        CheckStatus(v, &core.Audit{SubjectId: v.Model.ID})
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

func UpdateEntities(a *core.Audit, s core.Subject) {
    a.Status = AnalyseStatus(a, s)
    s.Status = a.Status

    core.App.DB.Save(a)
    core.App.DB.Save(&s)
}

// Checks the Status of a service/website using a HTTP "ping"
func CheckStatus(s core.Subject, a *core.Audit) {
    a.Result = false
    ts := time.Now()
    req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", s.Domain, s.PingURI), nil)
    req.SetBasicAuth(s.User, s.Pass)

    if err != nil {
        core.App.Log.Error(err)
        return
    }

    resp, err := HttpClient.Do(req)

    if err != nil {
        core.App.Log.Error(err)
        a.ResponseTime = time.Since(ts).Seconds()
        return
    }

    if resp.StatusCode == http.StatusOK {
        a.Result = true
        a.ResponseTime = time.Since(ts).Seconds()
    }

    if s.Advanced {
        CheckAdvancedStatus(s, a, resp)
    } else {
        CheckBasicStatus(s, a, resp)
    }
}

// Checks the basics, this is for endpoints that do not provide
// the detailed response specified in plugins/health.go
func CheckBasicStatus(s core.Subject, a *core.Audit, r *http.Response) {
    a.ResponseStatus = r.StatusCode
    UpdateEntities(a, s)
}

// This expects a more detailed response, check plugins/health.go for
// more information on what this expects
func CheckAdvancedStatus(s core.Subject, a *core.Audit, r *http.Response) {
    var h plugins.HealthCheckResponse

    if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
        core.App.Log.Errorf("Invalid json provided for %s - Checking using basic rules", s.Name)
        CheckBasicStatus(s, a, r)
        return
    }

    // Set vars
    s.Hostname = h.Host
    s.OS = h.OS
    s.Platform = h.Platform

    // Audit vars
    a.CPU = h.CPU
    a.Uptime = h.Uptime
    a.Memory = h.Memory
    a.ResponseStatus = r.StatusCode

    UpdateEntities(a, s)
}
