package site

import (
    "encoding/json"
    "github.com/Danzabar/WatchDog/core"
    "net/http"
)

// [GET|POST|PUT\PATCH] /ping
func Ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`Pong`))
}

// [GET] /status
func GetStatsPage(w http.ResponseWriter, r *http.Request) {

}

// [GET] /api/v1/subject
func GetSubjects(w http.ResponseWriter, r *http.Request) {
    var s []Subject

    if err := core.App.DB.Find(&s).Error; err != nil {
        core.WriteResponse(w, 500, core.RestResponse{Error: "Unable to fetch subjects"})
        return
    }

    js, _ := json.Marshal(&s)

    core.WriteResponseHeader(w, 200)
    w.Write(js)
}
