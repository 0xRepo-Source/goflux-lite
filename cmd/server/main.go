package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/0xRepo-Source/goflux-lite/pkg/auth"
	"github.com/0xRepo-Source/goflux-lite/pkg/config"
	"github.com/0xRepo-Source/goflux-lite/pkg/server"
	"github.com/0xRepo-Source/goflux-lite/pkg/storage"
)

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

func main() {
	configFile := flag.String("config", "goflux.json", "path to configuration file")
	port := flag.String("port", "", "server port (overrides config)")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux-lite-server version: 0.1.0-lite")
		return
	}

	// Load or create configuration
	cfg, err := config.LoadOrCreateConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override port if specified, or use internal IP for default config
	if *port != "" {
		internalIP := getInternalIP()
		cfg.Server.Address = fmt.Sprintf("%s:%s", internalIP, *port)
	} else if strings.Contains(cfg.Server.Address, "localhost") {
		// If config still has localhost, replace with internal IP
		internalIP := getInternalIP()
		parts := strings.Split(cfg.Server.Address, ":")
		if len(parts) == 2 {
			cfg.Server.Address = fmt.Sprintf("%s:%s", internalIP, parts[1])
		} else {
			cfg.Server.Address = fmt.Sprintf("%s:8080", internalIP)
		}
	}

	// Create storage backend
	store, err := storage.NewLocal(cfg.Server.StorageDir)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create server without web UI
	srv, err := server.New(store, cfg.Server.MetaDir)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Enable authentication if token file provided
	if cfg.Server.TokensFile != "" {
		tokenStore, err := auth.NewTokenStore(cfg.Server.TokensFile)
		if err != nil {
			log.Fatalf("Failed to load tokens: %v", err)
		}
		srv.EnableAuth(tokenStore)
		fmt.Printf("Authentication enabled: %s\n", cfg.Server.TokensFile)
	}

	// Create server config for sharing with clients
	serverConfig := &server.ServerConfig{
		Version:     "0.1.0-lite",
		AuthEnabled: cfg.Server.TokensFile != "",
	}
	serverConfig.Server.Address = cfg.Server.Address
	serverConfig.Server.StorageDir = cfg.Server.StorageDir
	serverConfig.Server.MetaDir = cfg.Server.MetaDir
	serverConfig.Server.TokensFile = cfg.Server.TokensFile
	serverConfig.Server.MaxFileSize = 1024 * 1024 * 1024 // 1GB default
	srv.SetConfig(serverConfig)

	// Enable discovery service
	if err := srv.EnableDiscovery(cfg.Server.Address, "0.1.0-lite"); err != nil {
		fmt.Printf("Warning: Failed to enable discovery: %v\n", err)
	}

	// Enable automatic firewall configuration
	srv.EnableFirewall(cfg.Server.Address)

	fmt.Printf("Starting goflux-lite server on %s\n", cfg.Server.Address)
	fmt.Printf("Storage directory: %s\n", cfg.Server.StorageDir)
	fmt.Printf("Configuration: %s\n", *configFile)

	// Start server
	if err := srv.Start(cfg.Server.Address); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
