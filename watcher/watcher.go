package watcher

import (
    "github.com/Danzabar/WatchDog/core"
    "net/http"
)

func Watch() {
    core.App.Log.Debug("Watching!")
}

// Interface for "Probes" this
// could be a SOAP client or HTTP Client etc
type Probe interface {
    CheckStatus(s core.Subject, a core.Audit) core.Audit
}

// HTTP Probe struct
type HttpProbe struct{}

// Checks the Status of a service/website using a HTTP "ping"
func (h *HttpProbe) CheckStatus(s core.Subject, a core.Audit) core.Audit {
    a.Status = false
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", s.Endpoint, s.PingURI), nil)

    if err != nil {
        core.App.Log.Error(err)
        return a
    }

    resp, err := http.DefaultClient.Do(req)

    if resp.StatusCode == http.StatusOK {
        a.Status = true

    }

    return a
}
