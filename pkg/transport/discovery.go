package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

// DiscoveredServer represents a server found on the network
type DiscoveredServer struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Address     string `json:"address"`
	Port        string `json:"port"`
	AuthEnabled bool   `json:"auth_enabled"`
	Timestamp   int64  `json:"timestamp"`
	LastSeen    time.Time
}

// Discovery client for finding GoFlux servers on the network
type DiscoveryClient struct {
	discovered map[string]*DiscoveredServer
	stopChan   chan struct{}
}

const (
	ClientDiscoveryPort    = 8081
	DiscoveryTimeout       = 5 * time.Second
	ServerExpiry           = 60 * time.Second
	DiscoveryMagicResponse = "GOFLUX-LITE-DISCOVERY"
)

// NewDiscoveryClient creates a new discovery client
func NewDiscoveryClient() *DiscoveryClient {
	return &DiscoveryClient{
		discovered: make(map[string]*DiscoveredServer),
		stopChan:   make(chan struct{}),
	}
}

// DiscoverServers scans the network for GoFlux servers
func (d *DiscoveryClient) DiscoverServers() ([]*DiscoveredServer, error) {
	// Listen for UDP broadcasts
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", ClientDiscoveryPort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP listener: %w", err)
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(DiscoveryTimeout))

	// Collect responses
	buffer := make([]byte, 1024)
	now := time.Now()

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// Timeout or other error, stop collecting
			break
		}

		// Parse message
		var message map[string]interface{}
		if err := json.Unmarshal(buffer[:n], &message); err != nil {
			continue // Invalid JSON, skip
		}

		// Check magic string
		magic, ok := message["magic"].(string)
		if !ok || magic != DiscoveryMagicResponse {
			continue // Not a GoFlux discovery message
		}

		// Extract server info
		dataInterface, ok := message["data"]
		if !ok {
			continue
		}

		dataBytes, err := json.Marshal(dataInterface)
		if err != nil {
			continue
		}

		var serverInfo DiscoveredServer
		if err := json.Unmarshal(dataBytes, &serverInfo); err != nil {
			continue
		}

		// Set discovery time and source
		serverInfo.LastSeen = now

		// Use the actual responding IP if address seems to be localhost/internal
		if strings.HasPrefix(serverInfo.Address, "localhost") ||
			strings.HasPrefix(serverInfo.Address, "127.0.0.1") ||
			strings.HasPrefix(serverInfo.Address, "0.0.0.0") {
			serverInfo.Address = fmt.Sprintf("%s:%s", remoteAddr.IP.String(), serverInfo.Port)
		}

		// Store unique server (by address)
		d.discovered[serverInfo.Address] = &serverInfo

		// Reset timeout to continue collecting
		conn.SetReadDeadline(time.Now().Add(time.Second))
	}

	// Clean up expired entries
	d.cleanupExpired()

	// Convert to slice and sort by last seen (newest first)
	var servers []*DiscoveredServer
	for _, server := range d.discovered {
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].LastSeen.After(servers[j].LastSeen)
	})

	return servers, nil
}

// cleanupExpired removes servers that haven't been seen recently
func (d *DiscoveryClient) cleanupExpired() {
	cutoff := time.Now().Add(-ServerExpiry)
	for addr, server := range d.discovered {
		if server.LastSeen.Before(cutoff) {
			delete(d.discovered, addr)
		}
	}
}

// GetServerConfig retrieves configuration from a discovered server
func (d *DiscoveryClient) GetServerConfig(serverAddr string) (map[string]interface{}, error) {
	// Ensure http:// prefix
	if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
		serverAddr = "http://" + serverAddr
	}

	// Request config from server
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(serverAddr + "/config")
	if err != nil {
		return nil, fmt.Errorf("failed to get config from server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

// FormatServerList returns a human-readable list of discovered servers
func (d *DiscoveryClient) FormatServerList(servers []*DiscoveredServer) string {
	if len(servers) == 0 {
		return "No GoFlux servers found on the local network.\n"
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Found %d GoFlux server(s):\n\n", len(servers)))

	for i, server := range servers {
		authStatus := "No Auth"
		if server.AuthEnabled {
			authStatus = "Auth Required"
		}

		age := time.Since(server.LastSeen)
		ageStr := "just now"
		if age > time.Second {
			ageStr = fmt.Sprintf("%ds ago", int(age.Seconds()))
		}

		output.WriteString(fmt.Sprintf("%d. %s (v%s)\n", i+1, server.Name, server.Version))
		output.WriteString(fmt.Sprintf("   Address: %s\n", server.Address))
		output.WriteString(fmt.Sprintf("   Status:  %s\n", authStatus))
		output.WriteString(fmt.Sprintf("   Seen:    %s\n", ageStr))
		output.WriteString("\n")
	}

	output.WriteString("Use 'gfl config <address>' to configure your client for a server.\n")
	return output.String()
}
