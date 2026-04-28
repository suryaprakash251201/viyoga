package linux

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// LogEntry represents a single journal log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Hostname  string `json:"hostname"`
	Unit      string `json:"unit"`
	Message   string `json:"message"`
	Priority  string `json:"priority"`
	PID       string `json:"pid"`
}

// LogFilter configures log query parameters.
type LogFilter struct {
	Unit     string `json:"unit,omitempty"`
	Lines    int    `json:"lines,omitempty"`
	Priority string `json:"priority,omitempty"` // emerg,alert,crit,err,warning,notice,info,debug
	Since    string `json:"since,omitempty"`    // e.g. "1 hour ago", "today"
	Until    string `json:"until,omitempty"`
	Grep     string `json:"grep,omitempty"`
}

// LogManager handles journal log reading.
type LogManager struct{}

func NewLogManager() *LogManager {
	return &LogManager{}
}

// GetLogs retrieves journal logs based on the provided filter.
func (m *LogManager) GetLogs(ctx context.Context, filter LogFilter) ([]LogEntry, error) {
	args := []string{"--no-pager", "--output=short-iso"}

	if filter.Unit != "" {
		if !isValidServiceName(filter.Unit) {
			return nil, fmt.Errorf("invalid unit name: %s", filter.Unit)
		}
		args = append(args, "-u", filter.Unit)
	}

	if filter.Lines > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", filter.Lines))
	} else {
		args = append(args, "-n", "200")
	}

	if filter.Priority != "" {
		args = append(args, "-p", filter.Priority)
	}

	if filter.Since != "" {
		args = append(args, "--since", filter.Since)
	}

	if filter.Until != "" {
		args = append(args, "--until", filter.Until)
	}

	if filter.Grep != "" && len(filter.Grep) < 256 {
		args = append(args, "-g", filter.Grep)
	}

	cmd := exec.CommandContext(ctx, "journalctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("journalctl: %w", err)
	}

	return parseJournalOutput(string(out)), nil
}

// StreamLogs returns a channel of log entries for real-time follow.
func (m *LogManager) StreamLogs(ctx context.Context, unit string) (<-chan LogEntry, error) {
	args := []string{"-f", "--no-pager", "--output=short-iso", "-n", "0"}
	if unit != "" {
		if !isValidServiceName(unit) {
			return nil, fmt.Errorf("invalid unit name: %s", unit)
		}
		args = append(args, "-u", unit)
	}

	cmd := exec.CommandContext(ctx, "journalctl", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("journalctl pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("journalctl start: %w", err)
	}

	ch := make(chan LogEntry, 100)
	go func() {
		defer close(ch)
		defer cmd.Wait()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				entries := parseJournalOutput(scanner.Text())
				if len(entries) > 0 {
					ch <- entries[0]
				}
			}
		}
	}()

	return ch, nil
}

func parseJournalOutput(output string) []LogEntry {
	var entries []LogEntry
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		entry := parseJournalLine(line)
		if entry.Message != "" {
			entries = append(entries, entry)
		}
	}
	return entries
}

func parseJournalLine(line string) LogEntry {
	// Format: "2024-01-01T12:00:00+0000 hostname unit[pid]: message"
	entry := LogEntry{}

	// Try to parse ISO timestamp
	if len(line) > 25 {
		if ts, err := time.Parse("2006-01-02T15:04:05-0700", line[:25]); err == nil {
			entry.Timestamp = ts.Format(time.RFC3339)
			line = strings.TrimSpace(line[25:])
		} else if len(line) > 15 {
			// Try syslog format: "Jan  1 12:00:00"
			entry.Timestamp = line[:15]
			line = strings.TrimSpace(line[15:])
		}
	}

	// Split remaining: hostname unit[pid]: message
	parts := strings.SplitN(line, " ", 3)
	if len(parts) >= 1 {
		entry.Hostname = parts[0]
	}
	if len(parts) >= 2 {
		unitPid := parts[1]
		if idx := strings.Index(unitPid, "["); idx > 0 {
			entry.Unit = unitPid[:idx]
			entry.PID = strings.TrimRight(unitPid[idx+1:], "]:")
		} else {
			entry.Unit = strings.TrimRight(unitPid, ":")
		}
	}
	if len(parts) >= 3 {
		entry.Message = parts[2]
	}

	return entry
}
