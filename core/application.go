package core

import (
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/op/go-logging"
    "net/http"
    "os"
)

var App *Application

/**
 * The Application Struct
 */
type Application struct {
    DB     *gorm.DB
    Router *mux.Router
    Log    *logging.Logger
    port   string
}

func NewApp(port string, dbDriver string, dbCreds string) {
    db, _ := gorm.Open(dbDriver, dbCreds)
    SetLogging()

    App = &Application{
        DB:     db,
        Router: mux.NewRouter(),
        Log:    logging.MustGetLogger("scribe"),
        port:   port,
    }
}

func SetLogging() {
    f := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`)
    b := logging.NewLogBackend(os.Stderr, "", 0)

    bf := logging.NewBackendFormatter(b, f)
    logging.SetBackend(bf)
}

func (a *Application) Run() {
    http.Handle("/", a.Router)
    a.Log.Debugf("Running app on " + a.port)
    a.Log.Critical(http.ListenAndServe(a.port, nil))
}
