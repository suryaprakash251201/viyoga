package collector

import (
	"context"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryMetrics holds RAM and swap usage data.
type MemoryMetrics struct {
	// RAM
	TotalBytes     uint64  `json:"total_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	FreeBytes      uint64  `json:"free_bytes"`
	UsagePercent   float64 `json:"usage_percent"`
	CachedBytes    uint64  `json:"cached_bytes"`
	BuffersBytes   uint64  `json:"buffers_bytes"`

	// Swap
	SwapTotalBytes   uint64  `json:"swap_total_bytes"`
	SwapUsedBytes    uint64  `json:"swap_used_bytes"`
	SwapFreeBytes    uint64  `json:"swap_free_bytes"`
	SwapUsagePercent float64 `json:"swap_usage_percent"`
}

// MemoryCollector collects memory metrics using gopsutil.
type MemoryCollector struct{}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (c *MemoryCollector) Name() string {
	return "memory"
}

func (c *MemoryCollector) Collect(ctx context.Context) (interface{}, error) {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	s, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		// Swap info is non-critical
		s = &mem.SwapMemoryStat{}
	}

	return &MemoryMetrics{
		TotalBytes:       v.Total,
		UsedBytes:        v.Used,
		AvailableBytes:   v.Available,
		FreeBytes:        v.Free,
		UsagePercent:     v.UsedPercent,
		CachedBytes:      v.Cached,
		BuffersBytes:     v.Buffers,
		SwapTotalBytes:   s.Total,
		SwapUsedBytes:    s.Used,
		SwapFreeBytes:    s.Free,
		SwapUsagePercent: s.UsedPercent,
	}, nil
}
