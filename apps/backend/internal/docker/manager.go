package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

// Manager wraps the Docker SDK client.
type Manager struct {
	cli *client.Client
}

// ContainerSummary is a simplified container view for the API.
type ContainerSummary struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Created int64             `json:"created"`
	Ports   []PortMapping     `json:"ports"`
	Labels  map[string]string `json:"labels"`
}

// PortMapping represents a container port mapping.
type PortMapping struct {
	Private uint16 `json:"private"`
	Public  uint16 `json:"public"`
	Type    string `json:"type"`
	IP      string `json:"ip"`
}

// ImageSummary is a simplified image view.
type ImageSummary struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    int64    `json:"size"`
	Created int64    `json:"created"`
}

// ContainerStats holds real-time stats for a container.
type ContainerStats struct {
	CPUPercent float64 `json:"cpu_percent"`
	MemUsage   uint64  `json:"mem_usage"`
	MemLimit   uint64  `json:"mem_limit"`
	MemPercent float64 `json:"mem_percent"`
	NetRx      uint64  `json:"net_rx"`
	NetTx      uint64  `json:"net_tx"`
	BlockRead  uint64  `json:"block_read"`
	BlockWrite uint64  `json:"block_write"`
	PIDs       uint64  `json:"pids"`
}

// NewManager creates a Docker manager. Returns nil manager if Docker is unavailable.
func NewManager(socketPath string) *Manager {
	opts := []client.Opt{client.FromEnv, client.WithAPIVersionNegotiation()}
	if socketPath != "" {
		opts = append(opts, client.WithHost(socketPath))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		log.Warn().Err(err).Msg("Docker client unavailable")
		return nil
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := cli.Ping(ctx); err != nil {
		log.Warn().Err(err).Msg("Docker daemon not reachable")
		return nil
	}

	log.Info().Msg("Docker client connected")
	return &Manager{cli: cli}
}

// IsAvailable returns true if Docker is connected.
func (m *Manager) IsAvailable() bool {
	return m != nil && m.cli != nil
}

// ListContainers returns all containers (running and stopped).
func (m *Manager) ListContainers(ctx context.Context, all bool) ([]ContainerSummary, error) {
	if !m.IsAvailable() {
		return nil, fmt.Errorf("docker not available")
	}

	containers, err := m.cli.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("docker list: %w", err)
	}

	var result []ContainerSummary
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		var ports []PortMapping
		for _, p := range c.Ports {
			ports = append(ports, PortMapping{
				Private: p.PrivatePort,
				Public:  p.PublicPort,
				Type:    p.Type,
				IP:      p.IP,
			})
		}

		result = append(result, ContainerSummary{
			ID:      c.ID[:12],
			Name:    name,
			Image:   c.Image,
			State:   c.State,
			Status:  c.Status,
			Created: c.Created,
			Ports:   ports,
			Labels:  c.Labels,
		})
	}
	return result, nil
}

// StartContainer starts a stopped container.
func (m *Manager) StartContainer(ctx context.Context, id string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("docker not available")
	}
	return m.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a running container.
func (m *Manager) StopContainer(ctx context.Context, id string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("docker not available")
	}
	timeout := 10
	return m.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

// RestartContainer restarts a container.
func (m *Manager) RestartContainer(ctx context.Context, id string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("docker not available")
	}
	timeout := 10
	return m.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer removes a container (force).
func (m *Manager) RemoveContainer(ctx context.Context, id string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("docker not available")
	}
	return m.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true})
}

// GetContainerLogs retrieves container logs.
func (m *Manager) GetContainerLogs(ctx context.Context, id string, lines int) (string, error) {
	if !m.IsAvailable() {
		return "", fmt.Errorf("docker not available")
	}

	tail := fmt.Sprintf("%d", lines)
	if lines <= 0 {
		tail = "100"
	}

	reader, err := m.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	})
	if err != nil {
		return "", fmt.Errorf("docker logs: %w", err)
	}
	defer reader.Close()

	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ListImages returns all Docker images.
func (m *Manager) ListImages(ctx context.Context) ([]ImageSummary, error) {
	if !m.IsAvailable() {
		return nil, fmt.Errorf("docker not available")
	}

	images, err := m.cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("docker images: %w", err)
	}

	var result []ImageSummary
	for _, img := range images {
		result = append(result, ImageSummary{
			ID:      img.ID[7:19], // sha256:xxxx -> first 12 chars
			Tags:    img.RepoTags,
			Size:    img.Size,
			Created: img.Created,
		})
	}
	return result, nil
}

// RemoveImage removes a Docker image.
func (m *Manager) RemoveImage(ctx context.Context, id string) error {
	if !m.IsAvailable() {
		return fmt.Errorf("docker not available")
	}
	_, err := m.cli.ImageRemove(ctx, id, image.RemoveOptions{Force: true, PruneChildren: true})
	return err
}

// PruneContainers removes all stopped containers.
func (m *Manager) PruneContainers(ctx context.Context) (uint64, error) {
	if !m.IsAvailable() {
		return 0, fmt.Errorf("docker not available")
	}
	report, err := m.cli.ContainersPrune(ctx, filters.Args{})
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// PruneImages removes unused images.
func (m *Manager) PruneImages(ctx context.Context) (uint64, error) {
	if !m.IsAvailable() {
		return 0, fmt.Errorf("docker not available")
	}
	report, err := m.cli.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

// InspectContainer returns detailed container info.
func (m *Manager) InspectContainer(ctx context.Context, id string) (types.ContainerJSON, error) {
	if !m.IsAvailable() {
		return types.ContainerJSON{}, fmt.Errorf("docker not available")
	}
	return m.cli.ContainerInspect(ctx, id)
}
