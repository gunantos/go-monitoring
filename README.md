# Go Monitoring

Monitoring server ringan berbasis Go untuk memantau status server **database** atau **app** melalui WebSocket.  
Mendukung multi-platform: Windows, Linux, macOS, dan ARM. Dilengkapi fitur **self-update otomatis** dari GitHub.

---

## Fitur

- Monitor CPU, RAM, dan load server.
- WebSocket server untuk push status secara real-time.
- Multi-platform: Windows, Linux, macOS, ARM64, ARM32.
- CLI fleksibel: pilih tipe server, port, dan interval update.
- Self-update otomatis dari GitHub release.
- Ringan, stabil, dan cross-platform.

---

## Instalasi

### Prasyarat

- Go >= 1.23
- Git (untuk clone repository)
- Sistem operasi Windows, Linux, atau macOS

### Clone repository

```bash
git clone https://github.com/gunantos/go-monitoring.git
cd go-monitoring
````

### Install dependencies

```bash
go mod tidy
```

### Build

#### Windows

```bash
GOOS=windows GOARCH=amd64 go build -o monitoring.exe main.go
```

#### Linux

```bash
GOOS=linux GOARCH=amd64 go build -o monitoring-linux main.go
```

#### macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o monitoring-macos main.go
```

---

## Penggunaan

```bash
# Jalankan monitoring server
./monitoring-linux -server=database -port=9800 -interval=2
```

### CLI Options

| Flag        | Default  | Deskripsi                          |
| ----------- | -------- | ---------------------------------- |
| `-server`   | database | Tipe server: `database` atau `app` |
| `-port`     | 9800     | Port WebSocket server              |
| `-interval` | 2        | Interval update status dalam detik |

---

## Self-Update

Server akan otomatis mengecek release terbaru di GitHub setiap X menit.
Jika ada versi baru, server akan memberitahu dan merekomendasikan restart.

---

## Integrasi Frontend

Koneksi WebSocket:

```js
const socket = io("ws://IP_SERVER:PORT/ws");

socket.on("serverStatusUpdate", (data) => {
    console.log(data);
});
```

Data dikirim dalam format:

```json
{
    "ip": "192.168.1.10",
    "cpuUsage": 12.5,
    "ramUsage": 45.2,
    "load1": 0.12,
    "load5": 0.10,
    "load15": 0.08,
    "serverType": "database"
}
```

---

## Build & Release Otomatis

Workflow GitHub Actions akan:

1. Membangun binary untuk multi-platform.
2. Upload artifact hasil build.
3. Membuat release GitHub otomatis ketika push tag `v*.*.*`.

---

## Contributing

* Fork repository
* Buat branch baru
* Commit perubahan
* Push ke branch
* Buat pull request

---

## Lisensi

MIT License © Gunanto

