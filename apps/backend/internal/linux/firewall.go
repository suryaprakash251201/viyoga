package linux

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// FirewallRule represents a UFW firewall rule.
type FirewallRule struct {
	Number    int    `json:"number"`
	To        string `json:"to"`
	Action    string `json:"action"`
	From      string `json:"from"`
	Comment   string `json:"comment,omitempty"`
	Direction string `json:"direction"` // "IN" or "OUT"
	V6        bool   `json:"v6"`
}

// FirewallStatus represents the overall firewall state.
type FirewallStatus struct {
	Active  bool           `json:"active"`
	Default string         `json:"default_policy"`
	Rules   []FirewallRule `json:"rules"`
}

// FirewallManager manages UFW firewall rules.
type FirewallManager struct{}

func NewFirewallManager() *FirewallManager {
	return &FirewallManager{}
}

// GetStatus returns the current UFW status and rules.
func (m *FirewallManager) GetStatus(ctx context.Context) (*FirewallStatus, error) {
	cmd := exec.CommandContext(ctx, "ufw", "status", "numbered")
	out, err := cmd.Output()
	if err != nil {
		return &FirewallStatus{Active: false}, nil
	}

	output := string(out)
	status := &FirewallStatus{
		Active: strings.Contains(output, "Status: active"),
	}

	// Parse rules
	ruleRegex := regexp.MustCompile(`\[\s*(\d+)\]\s+(.+?)\s+(ALLOW|DENY|REJECT|LIMIT)\s+(IN|OUT)?\s*(.*)`)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		matches := ruleRegex.FindStringSubmatch(line)
		if len(matches) >= 4 {
			num := 0
			fmt.Sscanf(matches[1], "%d", &num)
			rule := FirewallRule{
				Number:    num,
				To:        strings.TrimSpace(matches[2]),
				Action:    matches[3],
				Direction: matches[4],
				From:      strings.TrimSpace(matches[5]),
				V6:        strings.Contains(line, "(v6)"),
			}
			status.Rules = append(status.Rules, rule)
		}
	}

	return status, nil
}

// AllowPort adds an allow rule for a port.
func (m *FirewallManager) AllowPort(ctx context.Context, port string, proto string) error {
	return m.addRule(ctx, "allow", port, proto)
}

// DenyPort adds a deny rule for a port.
func (m *FirewallManager) DenyPort(ctx context.Context, port string, proto string) error {
	return m.addRule(ctx, "deny", port, proto)
}

// DeleteRule removes a rule by number.
func (m *FirewallManager) DeleteRule(ctx context.Context, ruleNumber int) error {
	if ruleNumber < 1 {
		return fmt.Errorf("invalid rule number: %d", ruleNumber)
	}
	cmd := exec.CommandContext(ctx, "ufw", "--force", "delete", fmt.Sprintf("%d", ruleNumber))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ufw delete: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// Enable enables UFW.
func (m *FirewallManager) Enable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ufw", "--force", "enable")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ufw enable: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// Disable disables UFW.
func (m *FirewallManager) Disable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ufw", "--force", "disable")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ufw disable: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func (m *FirewallManager) addRule(ctx context.Context, action, port, proto string) error {
	if !isValidPort(port) {
		return fmt.Errorf("invalid port: %s", port)
	}
	args := []string{action}
	if proto != "" && (proto == "tcp" || proto == "udp") {
		args = append(args, port+"/"+proto)
	} else {
		args = append(args, port)
	}

	cmd := exec.CommandContext(ctx, "ufw", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ufw %s: %s", action, strings.TrimSpace(string(out)))
	}
	return nil
}

func isValidPort(port string) bool {
	validPort := regexp.MustCompile(`^[0-9]+(/?(tcp|udp))?$|^[0-9]+:[0-9]+(/?(tcp|udp))?$`)
	return validPort.MatchString(port) && len(port) < 32
}
