package websockets

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"strings"
	"time"
)

type Stats struct {
	CPUUsage       float64      `json:"cpuUsage"`
	Memory         MemoryStats  `json:"memory"`
	Disks          []DiskStats  `json:"disk"`
	LoadAverage    LoadAvgStats `json:"loadAverage"`
	Timestamp      int64        `json:"timestamp"`
	Uptime         uint64       `json:"uptime"`
	CPUTemperature float64      `json:"cpuTemperature"`
}

type DiskStats struct {
	Name           string `json:"name"`
	TotalSpace     uint64 `json:"totalSpace"`
	AvailableSpace uint64 `json:"availableSpace"`
	UsedSpace      uint64 `json:"usedSpace"`
}

type MemoryStats struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

type LoadAvgStats struct {
	OneMinute      float64 `json:"oneMinute"`
	FiveMinutes    float64 `json:"fiveMinutes"`
	FifteenMinutes float64 `json:"fifteenMinutes"`
}

func StatsWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})

	// to check if client has closed its connection
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			// error thrown when client has disconnected so break the loop
			if err != nil {
				break
			}
		}
	}()

	for {
		select {
		// if go routine is closed, hence the client has disconnected so return and close the connection
		case <-done:
			return
		// calc stas and send to client every tick
		case <-ticker.C:
			stats, err := GetStats()
			if err != nil {
				continue
			}
			err = conn.WriteJSON(stats)
			// if there is some error sending the sats, return and close the connection
			if err != nil {
				return
			}
		}
	}
}

func GetStats() (*Stats, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskStats
	for _, p := range partitions {
		diskStat, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		disks = append(disks, DiskStats{
			Name:           p.Device,
			TotalSpace:     diskStat.Total,
			AvailableSpace: diskStat.Free,
			UsedSpace:      diskStat.Used,
		})
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, err
	}

	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}

	var cpuTemp float64
	var count int
	for _, t := range sensors {
		if strings.HasPrefix(t.SensorKey, "coretemp_packageid0") || strings.HasPrefix(t.SensorKey, "coretemp_core") {
			cpuTemp += t.Temperature
			count++
		}
	}
	if count > 0 {
		cpuTemp /= float64(count)
	}

	metrics := Stats{
		Timestamp: time.Now().Unix(),
		CPUUsage:  cpuPercent[0],
		Memory: MemoryStats{
			Total: memStat.Total,
			Used:  memStat.Used,
		},
		Disks: disks,
		LoadAverage: LoadAvgStats{
			OneMinute:      loadAvg.Load1,
			FiveMinutes:    loadAvg.Load5,
			FifteenMinutes: loadAvg.Load15,
		},
		Uptime:         uptime,
		CPUTemperature: cpuTemp,
	}
	fmt.Println("sent stats")
	return &metrics, nil
}
