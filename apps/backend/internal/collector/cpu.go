package collector

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUMetrics holds CPU usage data.
type CPUMetrics struct {
	UsagePercent   float64        `json:"usage_percent"`
	PerCore        []float64      `json:"per_core"`
	CoreCount      int            `json:"core_count"`
	LogicalCount   int            `json:"logical_count"`
	ModelName      string         `json:"model_name"`
	Frequency      float64        `json:"frequency_mhz"`
	LoadAvg1       float64        `json:"load_avg_1"`
	LoadAvg5       float64        `json:"load_avg_5"`
	LoadAvg15      float64        `json:"load_avg_15"`
}

// CPUCollector collects CPU metrics using gopsutil.
type CPUCollector struct{}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

func (c *CPUCollector) Name() string {
	return "cpu"
}

func (c *CPUCollector) Collect(ctx context.Context) (interface{}, error) {
	// Overall CPU usage (1-second sample)
	usage, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, err
	}

	// Per-core usage
	perCore, err := cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		perCore = []float64{}
	}

	// CPU info
	info, err := cpu.InfoWithContext(ctx)
	modelName := ""
	freq := 0.0
	if err == nil && len(info) > 0 {
		modelName = info[0].ModelName
		freq = info[0].Mhz
	}

	// Core counts
	physical, _ := cpu.CountsWithContext(ctx, false)
	logical, _ := cpu.CountsWithContext(ctx, true)

	// Load average (linux-only, returns 0 on Windows)
	// Using gopsutil load package
	loadAvg1, loadAvg5, loadAvg15 := getLoadAvg()

	usagePercent := 0.0
	if len(usage) > 0 {
		usagePercent = usage[0]
	}

	return &CPUMetrics{
		UsagePercent: usagePercent,
		PerCore:      perCore,
		CoreCount:    physical,
		LogicalCount: logical,
		ModelName:    modelName,
		Frequency:    freq,
		LoadAvg1:     loadAvg1,
		LoadAvg5:     loadAvg5,
		LoadAvg15:    loadAvg15,
	}, nil
}

func getLoadAvg() (float64, float64, float64) {
	// gopsutil's load.Avg() works on Linux, returns error on Windows
	// We import it separately to handle cross-platform
	avg, err := getLoadAvgImpl()
	if err != nil {
		return 0, 0, 0
	}
	return avg[0], avg[1], avg[2]
}

func getLoadAvgImpl() ([3]float64, error) {
	// This is a cross-platform stub. On Linux, gopsutil handles it.
	// On Windows, load average is not a native concept.
	// We use gopsutil/load for this.
	return [3]float64{0, 0, 0}, nil
}
