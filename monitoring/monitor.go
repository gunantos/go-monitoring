package monitoring

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type Monitor struct {
	ServerType string
	Port       int
	Interval   time.Duration
}

func NewMonitor(serverType string, port int, interval time.Duration) *Monitor {
	return &Monitor{
		ServerType: serverType,
		Port:       port,
		Interval:   interval,
	}
}

func (m *Monitor) Start() {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", m.Port),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, m.ServerType, m.Interval)
	})

	log.Printf("[%s] Monitoring server started on port %d", m.ServerType, m.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request, serverType string, interval time.Duration) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected:", r.RemoteAddr)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		status, err := getStatus(serverType)
		if err != nil {
			log.Println("Error getting status:", err)
			continue
		}
		if err := conn.WriteJSON(status); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func getStatus(serverType string) (*StatusData, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	ip := getLocalIP()

	return &StatusData{
		IP:       ip,
		CPUUsage: cpuPercents[0],
		RAMUsage: memStats.UsedPercent,
		Load1:    loadAvg.Load1,
		Load5:    loadAvg.Load5,
		Load15:   loadAvg.Load15,
		Server:   serverType,
	}, nil
}

// Get first non-loopback IPv4
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "0.0.0.0"
}
