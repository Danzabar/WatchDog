package main

import (
    "flag"
    "github.com/Danzabar/WatchDog/core"
    "github.com/Danzabar/WatchDog/site"
    "github.com/Danzabar/WatchDog/watcher"
    "github.com/jasonlvhit/gocron"
)

func main() {
    r := flag.Bool("r", false, "When included will run the application")
    m := flag.Bool("migrate", false, "Runs database schema migration tool if included")
    dd := flag.String("driver", "sqlite3", "The database driver to use")
    dc := flag.String("creds", "/tmp/main.db", "The database credentials")
    p := flag.String("port", ":8080", "The port on which this listens")
    w := flag.Bool("w", false, "Performs a watch operation on load if set")

    flag.Parse()

    core.NewApp(*p, *dd, *dc)

    site.Setup("site/templates/")

    if *m {
        Migrate()
    }

    if *w {
        watcher.Watch()
    }

    if *r {
        gocron.Every(1).Minute().Do(watcher.Watch)
        gocron.Start()
        core.App.Run()
    }
}

func Migrate() {
    core.App.Log.Debug("Starting Migrations")
    site.Migrate()
    core.App.Log.Debug("Finished Migrations")
}
