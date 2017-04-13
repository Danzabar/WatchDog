package site

import (
    "github.com/Danzabar/WatchDog/core"
)

func Setup() {
    core.App.Router.HandleFunc("/ping", Ping)
}
