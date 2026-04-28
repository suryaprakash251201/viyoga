package collector

import (
	"context"

	"github.com/shirou/gopsutil/v3/net"
)

// NetworkInterface holds stats for a single network interface.
type NetworkInterface struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	ErrIn       uint64 `json:"err_in"`
	ErrOut      uint64 `json:"err_out"`
	DropIn      uint64 `json:"drop_in"`
	DropOut     uint64 `json:"drop_out"`
}

// ConnectionSummary holds a summary of active connections.
type ConnectionSummary struct {
	TCP        int `json:"tcp"`
	UDP        int `json:"udp"`
	Listening  int `json:"listening"`
	Established int `json:"established"`
}

// NetworkMetrics holds all network-related data.
type NetworkMetrics struct {
	Interfaces  []NetworkInterface `json:"interfaces"`
	Connections ConnectionSummary  `json:"connections"`
	TotalSent   uint64             `json:"total_bytes_sent"`
	TotalRecv   uint64             `json:"total_bytes_recv"`
}

// NetworkCollector collects network interface metrics.
type NetworkCollector struct{}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{}
}

func (c *NetworkCollector) Name() string {
	return "network"
}

func (c *NetworkCollector) Collect(ctx context.Context) (interface{}, error) {
	counters, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	var ifaces []NetworkInterface
	var totalSent, totalRecv uint64

	for _, counter := range counters {
		ifaces = append(ifaces, NetworkInterface{
			Name:        counter.Name,
			BytesSent:   counter.BytesSent,
			BytesRecv:   counter.BytesRecv,
			PacketsSent: counter.PacketsSent,
			PacketsRecv: counter.PacketsRecv,
			ErrIn:       counter.Errin,
			ErrOut:      counter.Errout,
			DropIn:      counter.Dropin,
			DropOut:     counter.Dropout,
		})
		totalSent += counter.BytesSent
		totalRecv += counter.BytesRecv
	}

	// Connection summary
	connSummary := ConnectionSummary{}
	conns, err := net.ConnectionsWithContext(ctx, "all")
	if err == nil {
		for _, conn := range conns {
			switch conn.Type {
			case 1: // SOCK_STREAM (TCP)
				connSummary.TCP++
			case 2: // SOCK_DGRAM (UDP)
				connSummary.UDP++
			}
			switch conn.Status {
			case "LISTEN":
				connSummary.Listening++
			case "ESTABLISHED":
				connSummary.Established++
			}
		}
	}

	return &NetworkMetrics{
		Interfaces:  ifaces,
		Connections: connSummary,
		TotalSent:   totalSent,
		TotalRecv:   totalRecv,
	}, nil
}
