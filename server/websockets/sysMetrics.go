package websockets

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func GetMetrics() (map[string]interface{}, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"time":       time.Now(),
		"cpuPercent": cpuPercent,
		"memStat":    memStat,
		"diskStat":   diskStat,
		"loadAvg_1":  loadAvg.Load1,
		"loadAvg_5":  loadAvg.Load5,
		"loadAvg_15": loadAvg.Load15,
		"uptime":     uptime,
	}

	return metrics, nil
}
