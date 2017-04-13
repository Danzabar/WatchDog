package site

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/flosch/pongo2"
)

var (
    Template *pongo2.TemplateSet
)

func Setup() {
    fs := pongo2.MustNewLocalFileSystemLoader("templates/")
    Template = pongo2.NewSet("Templates", fs)

    // Ping route for AWS Cloudwatch
    core.App.Router.HandleFunc("/ping", Ping)

    // Route for stats/status page
    core.App.Router.HandleFunc("/status", GetStatsPage).Methods("GET")
}
