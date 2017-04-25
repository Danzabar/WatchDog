package site

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/flosch/pongo2"
    "net/http"
)

var (
    Template *pongo2.TemplateSet
)

func Setup(d string) {
    fs := pongo2.MustNewLocalFileSystemLoader(d)
    Template = pongo2.NewSet("Templates", fs)

    // Ping route
    core.App.Router.HandleFunc("/ping", Ping)

    // Route for stats/status page
    core.App.Router.HandleFunc("/", GetStatsPage).Methods("GET")

    // API
    core.App.Router.HandleFunc("/api/v1/subject", GetSubjects).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject/{id}", GetSubjectDetails).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject", PostSubject).Methods("POST")
}

func Migrate() {
    core.App.DB.AutoMigrate(&core.Audit{}, &core.Subject{})
}

// Func to Render a Pongo2 Template
func Render(n string, w http.ResponseWriter, ctx pongo2.Context) {
    tpl, e := Template.FromFile(n)

    if e != nil {
        core.App.Log.Error(e)
    }

    if err := tpl.ExecuteWriter(ctx, w); err != nil {
        core.App.Log.Error(err)
    }
}
