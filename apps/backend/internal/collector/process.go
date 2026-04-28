package collector

import (
	"context"
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo represents a running process.
type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	Username   string  `json:"username"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
	MemRSS     uint64  `json:"mem_rss"`
	Status     string  `json:"status"`
	CreateTime int64   `json:"create_time"`
	Cmdline    string  `json:"cmdline"`
	PPID       int32   `json:"ppid"`
	NumThreads int32   `json:"num_threads"`
}

// ProcessCollector collects process information.
type ProcessCollector struct{}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{}
}

func (c *ProcessCollector) Name() string { return "processes" }

func (c *ProcessCollector) Collect(ctx context.Context) (interface{}, error) {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	for _, p := range procs {
		info := ProcessInfo{PID: p.Pid}

		if name, err := p.NameWithContext(ctx); err == nil {
			info.Name = name
		}
		if user, err := p.UsernameWithContext(ctx); err == nil {
			info.Username = user
		}
		if cpu, err := p.CPUPercentWithContext(ctx); err == nil {
			info.CPUPercent = cpu
		}
		if mem, err := p.MemoryPercentWithContext(ctx); err == nil {
			info.MemPercent = mem
		}
		if memInfo, err := p.MemoryInfoWithContext(ctx); err == nil && memInfo != nil {
			info.MemRSS = memInfo.RSS
		}
		if status, err := p.StatusWithContext(ctx); err == nil && len(status) > 0 {
			info.Status = status[0]
		}
		if ct, err := p.CreateTimeWithContext(ctx); err == nil {
			info.CreateTime = ct
		}
		if cmd, err := p.CmdlineWithContext(ctx); err == nil {
			info.Cmdline = cmd
			if len(info.Cmdline) > 256 {
				info.Cmdline = info.Cmdline[:256]
			}
		}
		if ppid, err := p.PpidWithContext(ctx); err == nil {
			info.PPID = ppid
		}
		if threads, err := p.NumThreadsWithContext(ctx); err == nil {
			info.NumThreads = threads
		}

		processes = append(processes, info)
	}

	// Sort by CPU usage descending
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUPercent > processes[j].CPUPercent
	})

	// Limit to top 100
	if len(processes) > 100 {
		processes = processes[:100]
	}

	return processes, nil
}
