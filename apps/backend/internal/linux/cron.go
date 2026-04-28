package linux

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// CronJob represents a cron job entry.
type CronJob struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	User     string `json:"user"`
	Raw      string `json:"raw"`
}

// CronManager handles crontab operations.
type CronManager struct{}

func NewCronManager() *CronManager {
	return &CronManager{}
}

// ListCronJobs lists all cron jobs for the current user or a specified user.
func (m *CronManager) ListCronJobs(ctx context.Context, user string) ([]CronJob, error) {
	args := []string{"-l"}
	if user != "" {
		if !isValidServiceName(user) {
			return nil, fmt.Errorf("invalid username: %s", user)
		}
		args = append(args, "-u", user)
	}

	cmd := exec.CommandContext(ctx, "crontab", args...)
	out, err := cmd.Output()
	if err != nil {
		// "no crontab for user" is not a real error
		if strings.Contains(err.Error(), "no crontab") {
			return []CronJob{}, nil
		}
		return nil, fmt.Errorf("crontab: %w", err)
	}

	return parseCrontab(string(out), user), nil
}

func parseCrontab(output string, user string) []CronJob {
	var jobs []CronJob
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Cron format: min hour dom mon dow command
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		schedule := strings.Join(fields[:5], " ")
		command := strings.Join(fields[5:], " ")

		// Handle @reboot, @daily, etc.
		if strings.HasPrefix(line, "@") {
			schedule = fields[0]
			command = strings.Join(fields[1:], " ")
		}

		jobs = append(jobs, CronJob{
			Schedule: schedule,
			Command:  command,
			User:     user,
			Raw:      line,
		})
	}
	return jobs
}
