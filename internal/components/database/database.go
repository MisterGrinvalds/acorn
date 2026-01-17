// Package database provides database service management and status checking.
package database

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ServiceStatus represents the status of a database service.
type ServiceStatus struct {
	Name      string `json:"name" yaml:"name"`
	Installed bool   `json:"installed" yaml:"installed"`
	Running   bool   `json:"running" yaml:"running"`
	Status    string `json:"status" yaml:"status"`
}

// AllStatus contains status for all database services.
type AllStatus struct {
	Services []ServiceStatus `json:"services" yaml:"services"`
	Running  int             `json:"running" yaml:"running"`
	Stopped  int             `json:"stopped" yaml:"stopped"`
	Missing  int             `json:"missing" yaml:"missing"`
}

// Helper provides database management operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new database Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsDarwin returns true if running on macOS.
func (h *Helper) IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// CheckPostgreSQL checks PostgreSQL status.
func (h *Helper) CheckPostgreSQL() *ServiceStatus {
	status := &ServiceStatus{Name: "PostgreSQL"}

	// Check if pg_isready is available
	if _, err := exec.LookPath("pg_isready"); err != nil {
		if _, err := exec.LookPath("psql"); err != nil {
			status.Status = "Not installed"
			return status
		}
		status.Installed = true
		status.Status = "Installed (status unknown)"
		return status
	}

	status.Installed = true

	// Check if running
	cmd := exec.Command("pg_isready", "-q")
	if err := cmd.Run(); err == nil {
		status.Running = true
		status.Status = "Running"
	} else {
		status.Status = "Not running"
	}

	return status
}

// CheckMySQL checks MySQL status.
func (h *Helper) CheckMySQL() *ServiceStatus {
	status := &ServiceStatus{Name: "MySQL"}

	// Check if mysqladmin is available
	if _, err := exec.LookPath("mysqladmin"); err != nil {
		if _, err := exec.LookPath("mysql"); err != nil {
			status.Status = "Not installed"
			return status
		}
		status.Installed = true
		status.Status = "Installed (status unknown)"
		return status
	}

	status.Installed = true

	// Check if running
	cmd := exec.Command("mysqladmin", "ping", "-u", "root", "--silent")
	if err := cmd.Run(); err == nil {
		status.Running = true
		status.Status = "Running"
	} else {
		status.Status = "Not running"
	}

	return status
}

// CheckMongoDB checks MongoDB status.
func (h *Helper) CheckMongoDB() *ServiceStatus {
	status := &ServiceStatus{Name: "MongoDB"}

	// Check if mongosh is available
	if _, err := exec.LookPath("mongosh"); err != nil {
		status.Status = "Not installed"
		return status
	}

	status.Installed = true

	// Check if running
	cmd := exec.Command("mongosh", "--eval", "db.runCommand({ping:1})", "--quiet")
	out, err := cmd.Output()
	if err == nil && strings.Contains(string(out), "ok") {
		status.Running = true
		status.Status = "Running"
	} else {
		status.Status = "Not running"
	}

	return status
}

// CheckRedis checks Redis status.
func (h *Helper) CheckRedis() *ServiceStatus {
	status := &ServiceStatus{Name: "Redis"}

	// Check if redis-cli is available
	if _, err := exec.LookPath("redis-cli"); err != nil {
		status.Status = "Not installed"
		return status
	}

	status.Installed = true

	// Check if running
	cmd := exec.Command("redis-cli", "ping")
	out, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(out)) == "PONG" {
		status.Running = true
		status.Status = "Running"
	} else {
		status.Status = "Not running"
	}

	return status
}

// CheckNeo4j checks Neo4j status.
func (h *Helper) CheckNeo4j() *ServiceStatus {
	status := &ServiceStatus{Name: "Neo4j"}

	// Check if neo4j is available
	if _, err := exec.LookPath("neo4j"); err != nil {
		if _, err := exec.LookPath("cypher-shell"); err != nil {
			status.Status = "Not installed"
			return status
		}
		status.Installed = true
		status.Status = "Installed (status unknown)"
		return status
	}

	status.Installed = true

	// Check if running
	cmd := exec.Command("neo4j", "status")
	out, err := cmd.Output()
	if err == nil && strings.Contains(string(out), "running") {
		status.Running = true
		status.Status = "Running"
	} else {
		status.Status = "Not running"
	}

	return status
}

// CheckSQLite checks SQLite availability.
func (h *Helper) CheckSQLite() *ServiceStatus {
	status := &ServiceStatus{Name: "SQLite"}

	if _, err := exec.LookPath("sqlite3"); err != nil {
		status.Status = "Not installed"
		return status
	}

	status.Installed = true
	status.Running = true // SQLite is always "running" as it's file-based
	status.Status = "Available"

	return status
}

// GetAllStatus returns status of all database services.
func (h *Helper) GetAllStatus() *AllStatus {
	allStatus := &AllStatus{}

	services := []*ServiceStatus{
		h.CheckPostgreSQL(),
		h.CheckMySQL(),
		h.CheckMongoDB(),
		h.CheckRedis(),
		h.CheckNeo4j(),
		h.CheckSQLite(),
	}

	for _, svc := range services {
		allStatus.Services = append(allStatus.Services, *svc)
		if !svc.Installed {
			allStatus.Missing++
		} else if svc.Running {
			allStatus.Running++
		} else {
			allStatus.Stopped++
		}
	}

	return allStatus
}

// StartService starts a database service via brew (macOS only).
func (h *Helper) StartService(service string) error {
	if !h.IsDarwin() {
		return fmt.Errorf("service management is only supported on macOS with Homebrew")
	}

	brewService := h.getBrewServiceName(service)
	if brewService == "" {
		return fmt.Errorf("unknown service: %s", service)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: brew services start %s\n", brewService)
		return nil
	}

	cmd := exec.Command("brew", "services", "start", brewService)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StopService stops a database service via brew (macOS only).
func (h *Helper) StopService(service string) error {
	if !h.IsDarwin() {
		return fmt.Errorf("service management is only supported on macOS with Homebrew")
	}

	brewService := h.getBrewServiceName(service)
	if brewService == "" {
		return fmt.Errorf("unknown service: %s", service)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: brew services stop %s\n", brewService)
		return nil
	}

	cmd := exec.Command("brew", "services", "stop", brewService)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RestartService restarts a database service via brew (macOS only).
func (h *Helper) RestartService(service string) error {
	if !h.IsDarwin() {
		return fmt.Errorf("service management is only supported on macOS with Homebrew")
	}

	brewService := h.getBrewServiceName(service)
	if brewService == "" {
		return fmt.Errorf("unknown service: %s", service)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: brew services restart %s\n", brewService)
		return nil
	}

	cmd := exec.Command("brew", "services", "restart", brewService)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StartAll starts common database services (macOS only).
func (h *Helper) StartAll() error {
	if !h.IsDarwin() {
		return fmt.Errorf("service management is only supported on macOS with Homebrew")
	}

	services := []string{"postgresql@14", "redis", "mongodb-community"}

	for _, svc := range services {
		if h.dryRun {
			fmt.Printf("[dry-run] would run: brew services start %s\n", svc)
		} else {
			cmd := exec.Command("brew", "services", "start", svc)
			_ = cmd.Run() // Ignore errors for services that may not be installed
		}
	}

	return nil
}

// StopAll stops common database services (macOS only).
func (h *Helper) StopAll() error {
	if !h.IsDarwin() {
		return fmt.Errorf("service management is only supported on macOS with Homebrew")
	}

	services := []string{"postgresql@14", "redis", "mongodb-community"}

	for _, svc := range services {
		if h.dryRun {
			fmt.Printf("[dry-run] would run: brew services stop %s\n", svc)
		} else {
			cmd := exec.Command("brew", "services", "stop", svc)
			_ = cmd.Run() // Ignore errors for services that may not be installed
		}
	}

	return nil
}

// getBrewServiceName maps service names to brew service names.
func (h *Helper) getBrewServiceName(service string) string {
	switch strings.ToLower(service) {
	case "postgres", "postgresql", "pg":
		return "postgresql@14"
	case "mysql", "my":
		return "mysql"
	case "mongodb", "mongo":
		return "mongodb-community"
	case "redis", "rd":
		return "redis"
	case "neo4j", "neo":
		return "neo4j"
	case "kafka":
		return "kafka"
	case "zookeeper", "zk":
		return "zookeeper"
	default:
		return ""
	}
}

// GetSupportedServices returns list of supported database services.
func (h *Helper) GetSupportedServices() []string {
	return []string{
		"postgres (postgresql@14)",
		"mysql",
		"mongodb (mongodb-community)",
		"redis",
		"neo4j",
		"kafka",
		"zookeeper",
	}
}
