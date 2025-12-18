package monitoring

import (
	"time"

	"github.com/gorilla/websocket"
)

type StatusData struct {
	IP       string  `json:"ip"`
	CPUUsage float64 `json:"cpuUsage"`
	RAMUsage float64 `json:"ramUsage"`
	Load1    float64 `json:"load1"`
	Load5    float64 `json:"load5"`
	Load15   float64 `json:"load15"`
	Server   string  `json:"serverType"`
}

type Monitor struct {
	ServerType string
	Port       int
	Interval   time.Duration

	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}
