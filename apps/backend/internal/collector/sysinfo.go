package collector

import (
	"context"
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
)

// SystemInfoData holds static system information.
type SystemInfoData struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	PlatformFamily  string `json:"platform_family"`
	KernelVersion   string `json:"kernel_version"`
	KernelArch      string `json:"kernel_arch"`
	Uptime          uint64 `json:"uptime_seconds"`
	BootTime        uint64 `json:"boot_time"`
	Procs           uint64 `json:"procs"`
	GoVersion       string `json:"go_version"`
	ViyogaVersion   string `json:"viyoga_version"`
}

// SystemInfoCollector collects static system information.
type SystemInfoCollector struct {
	version string
}

func NewSystemInfoCollector(version string) *SystemInfoCollector {
	return &SystemInfoCollector{version: version}
}

func (c *SystemInfoCollector) Name() string {
	return "system_info"
}

func (c *SystemInfoCollector) Collect(ctx context.Context) (interface{}, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &SystemInfoData{
		Hostname:        info.Hostname,
		OS:              info.OS,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		PlatformFamily:  info.PlatformFamily,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		Uptime:          info.Uptime,
		BootTime:        info.BootTime,
		Procs:           info.Procs,
		GoVersion:       runtime.Version(),
		ViyogaVersion:   c.version,
	}, nil
}
