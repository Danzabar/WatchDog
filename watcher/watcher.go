package watcher

import (
    "fmt"
    "github.com/Danzabar/WatchDog/core"
    "net/http"
)

func Watch() {
    var s []core.Subject

    core.App.Log.Debug("Starting watcher...")
    core.App.DB.Find(&s)

    for _, v := range s {
        core.App.Log.Debugf("Checking %s", v.Domain)
        a := CheckStatus(v, &core.Audit{SubjectId: v.Model.ID})
        core.App.DB.Save(a)
        core.App.DB.Model(&v).Association("Audits").Append(a)
        core.App.Log.Debugf("Checked %s", v.Domain)
    }
}

// Checks the Status of a service/website using a HTTP "ping"
func CheckStatus(s core.Subject, a *core.Audit) *core.Audit {
    a.Result = false
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.Domain, s.PingURI), nil)

    if err != nil {
        core.App.Log.Error(err)
        return a
    }

    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        core.App.Log.Error(err)
        return a
    }

    if resp.StatusCode == http.StatusOK {
        a.Result = true
    }

    return a
}
