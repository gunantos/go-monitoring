package main

import (
	"time"

	"github.com/gunantos/go-monitoring/monitoring"
)

func main() {
	m := monitoring.NewMonitor("database", 9800, 2*time.Second)
	m.Start()
}
