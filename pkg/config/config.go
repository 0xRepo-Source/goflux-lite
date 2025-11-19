package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Address     string `json:"address"`     // Listen address (e.g., "0.0.0.0:80")
	StorageDir  string `json:"storage_dir"` // Storage directory path
	MetaDir     string `json:"meta_dir"`    // Metadata directory for resume
	TokensFile  string `json:"tokens_file"` // Path to tokens file (empty to disable auth)
	TLSCertFile string `json:"tls_cert"`    // TLS certificate file (empty for HTTP)
	TLSKeyFile  string `json:"tls_key"`     // TLS key file (empty for HTTP)
}

// ClientConfig holds client configuration
type ClientConfig struct {
	ServerURL string `json:"server_url"` // Server URL (e.g., "http://95.145.216.175")
	ChunkSize int    `json:"chunk_size"` // Chunk size in bytes
	Token     string `json:"token"`      // Authentication token (optional)
}

// Config holds both server and client configuration
type Config struct {
	Server ServerConfig `json:"server"`
	Client ClientConfig `json:"client"`
}

// getInternalIP returns the internal IP address of the machine
func getInternalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// Fallback: try to find first non-loopback interface
		interfaces, err := net.Interfaces()
		if err != nil {
			return "localhost"
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
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	internalIP := getInternalIP()
	return ServerConfig{
		Address:     fmt.Sprintf("%s:8080", internalIP),
		StorageDir:  "./data",
		MetaDir:     "./.goflux-meta",
		TokensFile:  "",
		TLSCertFile: "",
		TLSKeyFile:  "",
	}
}

// DefaultClientConfig returns default client configuration
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		ServerURL: "http://localhost",
		ChunkSize: 1024 * 1024, // 1MB
		Token:     "",
	}
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Server: DefaultServerConfig(),
		Client: DefaultClientConfig(),
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadOrCreateConfig loads config from file, or creates default if not exists
func LoadOrCreateConfig(path string) (*Config, error) {
	// Check if config exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config
		cfg := DefaultConfig()
		if err := SaveConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		fmt.Printf("Created default configuration at: %s\n", path)
		return &cfg, nil
	}

	// Load existing config
	return LoadConfig(path)
}
