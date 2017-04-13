package site

import (
    "github.com/Danzabar/WatchDog/core"
)

func Setup() {
    // Ping route for AWS Cloudwatch
    core.App.Router.HandleFunc("/ping", Ping)

    // Route for stats/status page
    core.App.Router.HandleFunc("/status", GetStatsPage).Methods("GET")
}
