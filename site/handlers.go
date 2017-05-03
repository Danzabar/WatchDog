package site

import (
    "github.com/flosch/pongo2"
    "net/http"
)

// [GET] /
func GetStatsPage(w http.ResponseWriter, r *http.Request) {
    Render("stats.html", w, pongo2.Context{})
}
