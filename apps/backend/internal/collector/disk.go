package collector

import (
	"context"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskPartition holds data about a single mount point.
type DiskPartition struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskIOStats holds I/O counters for a disk.
type DiskIOStats struct {
	Device     string `json:"device"`
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
	ReadTime   uint64 `json:"read_time_ms"`
	WriteTime  uint64 `json:"write_time_ms"`
}

// DiskMetrics holds all disk-related data.
type DiskMetrics struct {
	Partitions []DiskPartition `json:"partitions"`
	IO         []DiskIOStats   `json:"io"`
}

// DiskCollector collects disk space and I/O metrics.
type DiskCollector struct{}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{}
}

func (c *DiskCollector) Name() string {
	return "disk"
}

func (c *DiskCollector) Collect(ctx context.Context) (interface{}, error) {
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	var parts []DiskPartition
	for _, p := range partitions {
		usage, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil {
			continue
		}

		parts = append(parts, DiskPartition{
			Device:      p.Device,
			MountPoint:  p.Mountpoint,
			FSType:      p.Fstype,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	// Disk I/O counters
	ioCounters, err := disk.IOCountersWithContext(ctx)
	var ios []DiskIOStats
	if err == nil {
		for name, io := range ioCounters {
			ios = append(ios, DiskIOStats{
				Device:     name,
				ReadBytes:  io.ReadBytes,
				WriteBytes: io.WriteBytes,
				ReadCount:  io.ReadCount,
				WriteCount: io.WriteCount,
				ReadTime:   io.ReadTime,
				WriteTime:  io.WriteTime,
			})
		}
	}

	return &DiskMetrics{
		Partitions: parts,
		IO:         ios,
	}, nil
}
