package plugins

import (
    "encoding/json"
    "github.com/Danzabar/WatchDog/core"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "net/http"
)

// Response class for the health check endpoint
type HealthCheckResponse struct {
    Memory   float64 `json:"memory"`
    Uptime   uint64  `json:"uptime"`
    Host     string  `json:"host"`
    OS       string  `json:"os"`
    Platform string  `json:"platform"`
    CPU      float64 `json:"cpu"`
}

// Gathers a detailed list of information about the system
// this includes memory/cpu usage
// hostname, os, uptime
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // Memory usage
    v, _ := mem.VirtualMemory()
    u, _ := host.Info()
    c, _ := cpu.Percent(0, false)

    h := &HealthCheckResponse{
        Memory:   v.UsedPercent,
        Uptime:   u.Uptime,
        Host:     u.Hostname,
        OS:       u.OS,
        Platform: u.Platform,
        CPU:      c[0],
    }

    js, _ := json.Marshal(h)
    core.WriteResponseHeader(w, 200)
    w.Write(js)
}
