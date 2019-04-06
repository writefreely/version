package main

import (
	"fmt"
	"github.com/writefreely/version"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var (
	app *verApp
)

type verApp struct {
	// Config values
	Org   string
	Repo  string
	Port  string
	Debug bool

	Latest      string
	LastChecked *time.Time
	CacheTime   time.Duration
}

func main() {
	app = &verApp{}

	// Get config
	app.Org = os.Getenv("VER_ORG")
	app.Repo = os.Getenv("VER_REPO")
	app.Port = os.Getenv("VER_PORT")

	// Set some default values
	app.Debug = true
	app.CacheTime = 3 * time.Hour
	if app.Port == "" {
		app.Port = "8080"
	}

	// Validate config
	if app.Org == "" {
		log.Println("Missing GitHub organization variable: VER_ORG")
	}
	if app.Repo == "" {
		log.Println("Missing GitHub repository variable: VER_REPO")
	}
	if app.Org == "" || app.Repo == "" {
		os.Exit(1)
	}

	http.HandleFunc("/", handleGetCurrentVer)

	log.Printf("Serving on http://%s:%s\n", "localhost", app.Port)
	log.Printf("---")
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", "localhost", app.Port), nil)
	if err != nil {
		log.Printf("Unable to start: %v", err)
		os.Exit(1)
	}
}

func handleGetCurrentVer(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer logAndRecover(w, r, start)

	if app.LastChecked == nil || time.Since(*app.LastChecked) > app.CacheTime {
		if app.Debug {
			log.Printf("Last checked: %s; updating", app.LastChecked)
		}
		v, err := version.GetLatest(app.Org, app.Repo)
		if err != nil {
			log.Printf("Unable to get latest version: %v", err)
			if app.Latest == "" {
				log.Printf("No app version cached! Failing for user.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// Cache current value
			app.Latest = v
			now := time.Now()
			app.LastChecked = &now
		}
	} else {
		if app.Debug {
			log.Printf("Last checked: %s; NOT updating", app.LastChecked)
		}
	}

	fmt.Fprintf(w, app.Latest)
}

func logAndRecover(w http.ResponseWriter, r *http.Request, start time.Time) {
	if e := recover(); e != nil {
		log.Printf("%s: %s", e, debug.Stack())
		http.Error(w, "Error", http.StatusInternalServerError)
	}

	log.Printf(fmt.Sprintf("\"%s %s\" %s", r.Method, r.RequestURI, time.Since(start)))
}
