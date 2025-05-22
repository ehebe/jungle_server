package collector

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemStats struct {
	CPU float64 `json:"cpu"`
	Mem float64 `json:"mem"`
}

func Collect() *SystemStats {
	cpuPercents, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()

	return &SystemStats{
		CPU: cpuPercents[0],
		Mem: memStats.UsedPercent,
	}
}

func (s *SystemStats) ToJSON() []byte {
	data, _ := json.Marshal(s)
	return data
}
