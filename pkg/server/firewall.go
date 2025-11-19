package server

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// FirewallManager handles automatic firewall rule creation
type FirewallManager struct {
	serverPort    int
	discoveryPort int
}

// NewFirewallManager creates a new firewall manager
func NewFirewallManager(serverPort, discoveryPort int) *FirewallManager {
	return &FirewallManager{
		serverPort:    serverPort,
		discoveryPort: discoveryPort,
	}
}

// EnsureFirewallRules automatically creates firewall rules if needed
func (fm *FirewallManager) EnsureFirewallRules() {
	if runtime.GOOS != "windows" {
		// Only implement for Windows initially
		return
	}

	// Check if running as administrator
	if !fm.isAdmin() {
		fmt.Println("💡 For automatic firewall configuration, restart as Administrator")
		fmt.Println("   OR manually configure Windows Firewall:")
		fmt.Printf("   netsh advfirewall firewall add rule name=\"GoFlux Server\" dir=in action=allow protocol=TCP localport=%d\n", fm.serverPort)
		fmt.Printf("   netsh advfirewall firewall add rule name=\"GoFlux Discovery\" dir=in action=allow protocol=UDP localport=%d\n", fm.discoveryPort)
		fmt.Println()
		return
	}

	// Try to create firewall rules
	fmt.Println("🔥 Configuring Windows Firewall...")

	success := true

	// Create TCP rule for server
	if err := fm.createFirewallRule("GoFlux Server", "TCP", fm.serverPort); err != nil {
		fmt.Printf("⚠️  Failed to create server firewall rule: %v\n", err)
		success = false
	}

	// Create UDP rule for discovery
	if err := fm.createFirewallRule("GoFlux Discovery", "UDP", fm.discoveryPort); err != nil {
		fmt.Printf("⚠️  Failed to create discovery firewall rule: %v\n", err)
		success = false
	}

	if success {
		fmt.Println("✅ Firewall rules configured successfully")
	} else {
		fmt.Println("⚠️  Some firewall rules may need manual configuration")
	}
}

// isAdmin checks if the current process is running as administrator
func (fm *FirewallManager) isAdmin() bool {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

// createFirewallRule creates a Windows firewall rule
func (fm *FirewallManager) createFirewallRule(name, protocol string, port int) error {
	// Check if rule already exists
	checkCmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", fmt.Sprintf("name=%s", name))
	if err := checkCmd.Run(); err == nil {
		// Rule already exists
		fmt.Printf("   Firewall rule '%s' already exists\n", name)
		return nil
	}

	// Create the rule
	args := []string{
		"advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", name),
		"dir=in",
		"action=allow",
		fmt.Sprintf("protocol=%s", protocol),
		fmt.Sprintf("localport=%d", port),
	}

	cmd := exec.Command("netsh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("netsh failed: %v - %s", err, string(output))
	}

	fmt.Printf("   Created firewall rule: %s (%s:%d)\n", name, protocol, port)
	return nil
}

// RemoveFirewallRules removes the firewall rules (cleanup)
func (fm *FirewallManager) RemoveFirewallRules() {
	if runtime.GOOS != "windows" {
		return
	}

	if !fm.isAdmin() {
		return
	}

	// Remove rules (best effort, don't report errors)
	exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=GoFlux Server").Run()
	exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=GoFlux Discovery").Run()
}

// parsePortFromAddress extracts port number from address string
func parsePortFromAddress(address string) int {
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		return 8080 // default
	}

	port, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 8080 // default
	}

	return port
}
