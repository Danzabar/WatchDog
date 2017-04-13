package site

import (
    "net/http"
)

// [GET|POST|PUT\PATCH] /ping
func Ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`Pong`))
}

func GetStatsPage(w http.ResponseWriter, r *http.Request) {

}
