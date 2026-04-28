package linux

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// User represents a system user.
type User struct {
	Username string   `json:"username"`
	UID      int      `json:"uid"`
	GID      int      `json:"gid"`
	Comment  string   `json:"comment"`
	Home     string   `json:"home"`
	Shell    string   `json:"shell"`
	Groups   []string `json:"groups,omitempty"`
	IsSystem bool     `json:"is_system"`
}

// UserManager handles system user operations.
type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

// ListUsers returns all system users (optionally filtering non-system users).
func (m *UserManager) ListUsers(ctx context.Context, includeSystem bool) ([]User, error) {
	cmd := exec.CommandContext(ctx, "getent", "passwd")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getent passwd: %w", err)
	}

	var users []User
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}
		uid, _ := strconv.Atoi(parts[2])
		gid, _ := strconv.Atoi(parts[3])
		isSystem := uid < 1000 && uid != 0

		if !includeSystem && isSystem {
			continue
		}

		user := User{
			Username: parts[0],
			UID:      uid,
			GID:      gid,
			Comment:  parts[4],
			Home:     parts[5],
			Shell:    parts[6],
			IsSystem: isSystem,
		}

		// Get groups
		if groups, err := m.getUserGroups(ctx, parts[0]); err == nil {
			user.Groups = groups
		}

		users = append(users, user)
	}

	return users, nil
}

func (m *UserManager) getUserGroups(ctx context.Context, username string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "groups", username)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	output := strings.TrimSpace(string(out))
	// Output format: "username : group1 group2 group3"
	parts := strings.SplitN(output, ":", 2)
	if len(parts) < 2 {
		return nil, nil
	}
	groups := strings.Fields(strings.TrimSpace(parts[1]))
	return groups, nil
}
