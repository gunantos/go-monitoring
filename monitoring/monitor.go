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

func NewMonitor(serverType string, port int, interval time.Duration) *Monitor {
	return &Monitor{
		ServerType: serverType,
		Port:       port,
		Interval:   interval,
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (m *Monitor) Start() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		m.register <- conn
		go m.handleClient(conn)
	})

	go m.run()

	addr := fmt.Sprintf(":%d", m.Port)
	log.Printf("[%s] Monitoring server started on port %d", m.ServerType, m.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) run() {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()
	for {
		select {
		case conn := <-m.register:
			m.clients[conn] = true
			conn.WriteJSON(map[string]string{"event": "server_connect"})
			log.Println("Client connected:", conn.RemoteAddr())
		case conn := <-m.unregister:
			if _, ok := m.clients[conn]; ok {
				delete(m.clients, conn)
				conn.Close()
				log.Println("Client disconnected:", conn.RemoteAddr())
			}
		case <-ticker.C:
			status, err := getStatus(m.ServerType)
			if err != nil {
				log.Println("Error getting status:", err)
				continue
			}
			for client := range m.clients {
				err := client.WriteJSON(map[string]interface{}{
					"event":  "server_status",
					"status": status,
				})
				if err != nil {
					m.unregister <- client
				}
			}
		}
	}
}

func (m *Monitor) handleClient(conn *websocket.Conn) {
	defer func() { m.unregister <- conn }()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
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

	loadAvg, _ := load.Avg() // Windows load akan selalu 0, aman

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
