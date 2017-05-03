package site

import (
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/plugins"
    "github.com/flosch/pongo2"
    "net/http"
)

var (
    Template *pongo2.TemplateSet
)

func Setup(d string) {
    fs := pongo2.MustNewLocalFileSystemLoader(d)
    Template = pongo2.NewSet("Templates", fs)

    // Route for stats/status page
    core.App.Router.HandleFunc("/", GetStatsPage).Methods("GET")

    // HealthCheck
    core.App.Router.HandleFunc("/health", plugins.HealthCheckHandler).Methods("GET")

    // Assets
    core.App.Router.
        PathPrefix("/").
        Handler(http.StripPrefix("/", http.FileServer(http.Dir("site/assets/"))))
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
