package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Service represents a systemd service unit.
type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LoadState   string `json:"load_state"`
	ActiveState string `json:"active_state"`
	SubState    string `json:"sub_state"`
	UnitFile    string `json:"unit_file"`
	Enabled     string `json:"enabled"`
	Preset      string `json:"preset"`
}

// ServiceDetail contains extended information about a service.
type ServiceDetail struct {
	Service
	MainPID        int       `json:"main_pid"`
	MemoryCurrent  uint64    `json:"memory_current"`
	CPUUsage       string    `json:"cpu_usage"`
	ActiveEnter    time.Time `json:"active_enter_timestamp,omitempty"`
	ExecMainStart  time.Time `json:"exec_main_start,omitempty"`
	FragmentPath   string    `json:"fragment_path"`
	Restart        string    `json:"restart"`
	RestartUSec    string    `json:"restart_usec"`
	Type           string    `json:"type"`
	ActiveDuration string    `json:"active_duration"`
}

// SystemdManager handles systemd service operations.
type SystemdManager struct{}

func NewSystemdManager() *SystemdManager {
	return &SystemdManager{}
}

// ListServices returns all loaded systemd services.
func (m *SystemdManager) ListServices(ctx context.Context) ([]Service, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "list-units", "--type=service", "--all", "--no-pager", "--output=json")
	out, err := cmd.Output()
	if err != nil {
		// Fallback: parse text output on systems without JSON output
		return m.listServicesFallback(ctx)
	}

	var units []struct {
		Unit        string `json:"unit"`
		Load        string `json:"load"`
		Active      string `json:"active"`
		Sub         string `json:"sub"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(out, &units); err != nil {
		return m.listServicesFallback(ctx)
	}

	services := make([]Service, 0, len(units))
	for _, u := range units {
		if !strings.HasSuffix(u.Unit, ".service") {
			continue
		}
		services = append(services, Service{
			Name:        strings.TrimSuffix(u.Unit, ".service"),
			Description: u.Description,
			LoadState:   u.Load,
			ActiveState: u.Active,
			SubState:    u.Sub,
		})
	}
	return services, nil
}

func (m *SystemdManager) listServicesFallback(ctx context.Context) ([]Service, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl list-units: %w", err)
	}

	var services []Service
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		name := strings.TrimSuffix(fields[0], ".service")
		desc := ""
		if len(fields) > 4 {
			desc = strings.Join(fields[4:], " ")
		}
		services = append(services, Service{
			Name:        name,
			LoadState:   fields[1],
			ActiveState: fields[2],
			SubState:    fields[3],
			Description: desc,
		})
	}
	return services, nil
}

// GetService returns detailed information about a specific service.
func (m *SystemdManager) GetService(ctx context.Context, name string) (*ServiceDetail, error) {
	if !isValidServiceName(name) {
		return nil, fmt.Errorf("invalid service name: %s", name)
	}

	unit := name + ".service"
	props := []string{
		"Description", "LoadState", "ActiveState", "SubState",
		"MainPID", "MemoryCurrent", "CPUUsageNSec",
		"ActiveEnterTimestamp", "ExecMainStartTimestamp",
		"FragmentPath", "UnitFileState", "UnitFilePreset",
		"Restart", "RestartUSec", "Type",
	}

	cmd := exec.CommandContext(ctx, "systemctl", "show", unit, "--property="+strings.Join(props, ","))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl show %s: %w", name, err)
	}

	detail := &ServiceDetail{
		Service: Service{Name: name},
	}

	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		switch key {
		case "Description":
			detail.Description = val
		case "LoadState":
			detail.LoadState = val
		case "ActiveState":
			detail.ActiveState = val
		case "SubState":
			detail.SubState = val
		case "MainPID":
			fmt.Sscanf(val, "%d", &detail.MainPID)
		case "MemoryCurrent":
			if val != "[not set]" {
				fmt.Sscanf(val, "%d", &detail.MemoryCurrent)
			}
		case "CPUUsageNSec":
			if val != "[not set]" {
				detail.CPUUsage = val
			}
		case "FragmentPath":
			detail.FragmentPath = val
		case "UnitFileState":
			detail.Enabled = val
		case "UnitFilePreset":
			detail.Preset = val
		case "Restart":
			detail.Restart = val
		case "RestartUSec":
			detail.RestartUSec = val
		case "Type":
			detail.Type = val
		}
	}
	return detail, nil
}

// StartService starts a systemd service.
func (m *SystemdManager) StartService(ctx context.Context, name string) error {
	return m.runServiceCommand(ctx, "start", name)
}

// StopService stops a systemd service.
func (m *SystemdManager) StopService(ctx context.Context, name string) error {
	return m.runServiceCommand(ctx, "stop", name)
}

// RestartService restarts a systemd service.
func (m *SystemdManager) RestartService(ctx context.Context, name string) error {
	return m.runServiceCommand(ctx, "restart", name)
}

// EnableService enables a systemd service.
func (m *SystemdManager) EnableService(ctx context.Context, name string) error {
	return m.runServiceCommand(ctx, "enable", name)
}

// DisableService disables a systemd service.
func (m *SystemdManager) DisableService(ctx context.Context, name string) error {
	return m.runServiceCommand(ctx, "disable", name)
}

func (m *SystemdManager) runServiceCommand(ctx context.Context, action, name string) error {
	if !isValidServiceName(name) {
		return fmt.Errorf("invalid service name: %s", name)
	}
	unit := name + ".service"
	cmd := exec.CommandContext(ctx, "systemctl", action, unit)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl %s %s: %s", action, name, strings.TrimSpace(string(out)))
	}
	return nil
}

// isValidServiceName prevents command injection via service names.
func isValidServiceName(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '@' || c == '.') {
			return false
		}
	}
	return len(name) > 0 && len(name) < 256
}
