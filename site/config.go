package site

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/flosch/pongo2"
)

var (
    Template *pongo2.TemplateSet
)

func Setup(d string) {
    fs := pongo2.MustNewLocalFileSystemLoader(d)
    Template = pongo2.NewSet("Templates", fs)

    // Ping route for AWS Cloudwatch
    core.App.Router.HandleFunc("/ping", Ping)

    // Route for stats/status page
    core.App.Router.HandleFunc("/status", GetStatsPage).Methods("GET")

    // API
    core.App.Router.HandleFunc("/api/v1/subject", GetSubjects).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject/{id}", GetSubjectDetails).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject", PostSubject).Methods("POST")
}

func Migrate() {
    core.App.DB.AutoMigrate(&core.Audit{}, &core.Subject{})
}
