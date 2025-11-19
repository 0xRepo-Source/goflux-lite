package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// DiscoveryInfo represents server information broadcast on the network
type DiscoveryInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Address     string `json:"address"`
	Port        string `json:"port"`
	AuthEnabled bool   `json:"auth_enabled"`
	Timestamp   int64  `json:"timestamp"`
}

// DiscoveryService handles UDP broadcast announcements
type DiscoveryService struct {
	info     DiscoveryInfo
	conn     *net.UDPConn
	stopChan chan struct{}
}

const (
	DiscoveryPort     = 8081
	BroadcastInterval = 30 * time.Second
	DiscoveryMagic    = "GOFLUX-LITE-DISCOVERY"
)

// NewDiscoveryService creates a new discovery service
func NewDiscoveryService(serverAddress, version string, authEnabled bool) (*DiscoveryService, error) {
	// Parse server address to get port
	parts := strings.Split(serverAddress, ":")
	var port string
	if len(parts) == 2 {
		port = parts[1]
	} else {
		port = "8080" // default
	}

	info := DiscoveryInfo{
		Name:        "GoFlux Lite Server",
		Version:     version,
		Address:     serverAddress,
		Port:        port,
		AuthEnabled: authEnabled,
	}

	// Create UDP connection for broadcasting
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", DiscoveryPort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	return &DiscoveryService{
		info:     info,
		conn:     conn,
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins broadcasting server information
func (d *DiscoveryService) Start() {
	go d.broadcastLoop()
	fmt.Printf("Discovery service started on UDP port %d\n", DiscoveryPort)
}

// Stop halts the discovery service
func (d *DiscoveryService) Stop() {
	close(d.stopChan)
	if d.conn != nil {
		d.conn.Close()
	}
}

// broadcastLoop continuously broadcasts server information
func (d *DiscoveryService) broadcastLoop() {
	ticker := time.NewTicker(BroadcastInterval)
	defer ticker.Stop()

	// Send initial broadcast
	d.broadcast()

	for {
		select {
		case <-ticker.C:
			d.broadcast()
		case <-d.stopChan:
			return
		}
	}
}

// broadcast sends server information to the network
func (d *DiscoveryService) broadcast() {
	// Update timestamp
	d.info.Timestamp = time.Now().Unix()

	// Create broadcast message
	message := map[string]interface{}{
		"magic": DiscoveryMagic,
		"data":  d.info,
	}

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Failed to marshal discovery data: %v\n", err)
		return
	}

	// Get broadcast addresses
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Failed to get network interfaces: %v\n", err)
		return
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				// Calculate broadcast address
				broadcast := make(net.IP, 4)
				for i := 0; i < 4; i++ {
					broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
				}

				// Send broadcast
				broadcastAddr := &net.UDPAddr{
					IP:   broadcast,
					Port: DiscoveryPort,
				}

				_, err := d.conn.WriteToUDP(data, broadcastAddr)
				if err != nil {
					// Don't log every broadcast error to avoid spam
					continue
				}
			}
		}
	}
}
