package core

import (
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/op/go-logging"
    "gopkg.in/go-playground/validator.v9"
    "net/http"
    "os"
)

var App *Application

/**
 * The Application Struct
 */
type Application struct {
    DB        *gorm.DB
    Router    *mux.Router
    Log       *logging.Logger
    Alerts    bool
    Port      string
    Validator *validator.Validate
}

func NewApp(port string, dbDriver string, dbCreds string, alerts bool) {
    db, _ := gorm.Open(dbDriver, dbCreds)
    SetLogging()

    App = &Application{
        DB:        db,
        Router:    mux.NewRouter(),
        Log:       logging.MustGetLogger("scribe"),
        Port:      port,
        Alerts:    alerts,
        Validator: validator.New(),
    }
}

/**
 * Sets up logging for the application
 */
func SetLogging() {
    f := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s}%{color:reset} %{message}`)
    b := logging.NewLogBackend(os.Stderr, "", 0)

    bf := logging.NewBackendFormatter(b, f)
    logging.SetBackend(bf)
}

func (a *Application) Run() {
    http.Handle("/", a.Router)
    a.Log.Debugf("Running app on " + a.Port)
    a.Log.Critical(http.ListenAndServe(a.Port, nil))
}
