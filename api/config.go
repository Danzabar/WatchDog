package api

import (
    "github.com/Danzabar/WatchDog/core"
)

func Setup() {

    // API
    core.App.Router.HandleFunc("/api/v1/subject", GetSubjects).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject/{id}", GetSubjectDetails).Methods("GET")
    core.App.Router.HandleFunc("/api/v1/subject", PostSubject).Methods("POST")
    core.App.Router.HandleFunc("/api/v1/subject/{id}", DeleteSubject).Methods("DELETE")
}
