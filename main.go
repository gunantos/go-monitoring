package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blang/semver"
	"github.com/gunantos/go-monitoring/monitoring"
	selfupdate "github.com/rhysd/go-github-selfupdate/selfupdate"
)

var currentVersion = semver.MustParse("1.0.0")

func main() {
	// CLI flags
	serverType := flag.String("server", "database", "Server type: database or app")
	port := flag.Int("port", 9800, "Port to run the monitoring server on")
	interval := flag.Int("interval", 2, "Status update interval in seconds")
	updateCheck := flag.Int("update", 10, "Self-update check interval in minutes")
	flag.Parse()

	fmt.Printf("[%s] Starting monitoring server on port %d (version %s)\n", *serverType, *port, currentVersion)
	go autoUpdate(*updateCheck)
	m := monitoring.NewMonitor(*serverType, *port, time.Duration(*interval)*time.Second)
	m.Start()
}

func autoUpdate(minutes int) {
	for {
		time.Sleep(time.Duration(minutes) * time.Minute)
		latest, err := selfupdate.UpdateSelf(currentVersion, "gunantos/go-monitoring")
		if err != nil {
			log.Println("Self-update error:", err)
			continue
		}

		if !latest.Version.Equals(currentVersion) {
			log.Printf("Updated to version %s\n", latest.Version)
			log.Println("Restart server to apply the update...")
			os.Exit(0)
		} else {
			log.Println("No new version available")
		}
	}
}
