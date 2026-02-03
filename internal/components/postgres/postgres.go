// Package postgres provides PostgreSQL helper functionality.
package postgres

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents PostgreSQL installation status.
type Status struct {
	Installed       bool   `json:"installed" yaml:"installed"`
	Version         string `json:"version,omitempty" yaml:"version,omitempty"`
	PsqlVersion     string `json:"psql_version,omitempty" yaml:"psql_version,omitempty"`
	ServerRunning   bool   `json:"server_running" yaml:"server_running"`
	DockerRunning   bool   `json:"docker_running" yaml:"docker_running"`
	ContainerID     string `json:"container_id,omitempty" yaml:"container_id,omitempty"`
	DataDir         string `json:"data_dir,omitempty" yaml:"data_dir,omitempty"`
	Port            string `json:"port,omitempty" yaml:"port,omitempty"`
}

// Database represents a PostgreSQL database.
type Database struct {
	Name  string `json:"name" yaml:"name"`
	Owner string `json:"owner,omitempty" yaml:"owner,omitempty"`
	Size  string `json:"size,omitempty" yaml:"size,omitempty"`
}

// Connection represents a database connection config.
type Connection struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
	User     string `json:"user" yaml:"user"`
}

// Helper provides PostgreSQL helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new PostgreSQL Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns PostgreSQL status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		Port: "5432",
	}

	// Check psql
	out, err := exec.Command("psql", "--version").Output()
	if err == nil {
		status.Installed = true
		version := strings.TrimSpace(string(out))
		// Extract version number
		parts := strings.Fields(version)
		if len(parts) >= 3 {
			status.PsqlVersion = parts[2]
		} else {
			status.PsqlVersion = version
		}
	}

	// Check postgres server
	out, err = exec.Command("postgres", "--version").Output()
	if err == nil {
		version := strings.TrimSpace(string(out))
		parts := strings.Fields(version)
		if len(parts) >= 3 {
			status.Version = parts[2]
		}
	}

	// Check if server is running (Homebrew services)
	out, _ = exec.Command("brew", "services", "info", "postgresql", "--json").Output()
	if strings.Contains(string(out), `"running":true`) {
		status.ServerRunning = true
	}

	// Check for pg_isready
	if err := exec.Command("pg_isready", "-q").Run(); err == nil {
		status.ServerRunning = true
	}

	// Check Docker container
	containerID := h.getDockerContainer()
	if containerID != "" {
		status.DockerRunning = true
		status.ContainerID = containerID
		status.Port = h.getContainerPort(containerID)
	}

	// Get data directory
	status.DataDir = h.getDataDir()

	return status
}

// getDockerContainer returns running PostgreSQL container ID.
func (h *Helper) getDockerContainer() string {
	cmd := exec.Command("docker", "ps", "-q", "--filter", "ancestor=postgres")
	out, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return strings.TrimSpace(string(out))
	}

	// Try by name
	cmd = exec.Command("docker", "ps", "-q", "--filter", "name=postgres")
	out, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return ""
}

// getContainerPort returns the mapped port for a container.
func (h *Helper) getContainerPort(containerID string) string {
	cmd := exec.Command("docker", "port", containerID, "5432")
	out, err := cmd.Output()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(out)), ":")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}
	return "5432"
}

// getDataDir returns the PostgreSQL data directory.
func (h *Helper) getDataDir() string {
	// Check PGDATA environment variable
	if pgdata := os.Getenv("PGDATA"); pgdata != "" {
		return pgdata
	}

	// Common Homebrew location
	home, _ := os.UserHomeDir()
	homebrewData := filepath.Join(home, "Library/Application Support/Postgres")
	if _, err := os.Stat(homebrewData); err == nil {
		return homebrewData
	}

	// Check /usr/local/var/postgres
	if _, err := os.Stat("/usr/local/var/postgres"); err == nil {
		return "/usr/local/var/postgres"
	}

	// Homebrew on Apple Silicon
	if _, err := os.Stat("/opt/homebrew/var/postgres"); err == nil {
		return "/opt/homebrew/var/postgres"
	}

	return ""
}

// Start starts PostgreSQL server.
func (h *Helper) Start(useDocker bool, port int, password string) error {
	if port == 0 {
		port = 5432
	}
	if password == "" {
		password = "postgres"
	}

	if h.dryRun {
		if useDocker {
			fmt.Printf("[dry-run] would start PostgreSQL via Docker on port %d\n", port)
		} else {
			fmt.Println("[dry-run] would start PostgreSQL via Homebrew")
		}
		return nil
	}

	if useDocker {
		return h.startDocker(port, password)
	}

	return h.startLocal()
}

// startLocal starts PostgreSQL via Homebrew.
func (h *Helper) startLocal() error {
	fmt.Println("Starting PostgreSQL...")
	cmd := exec.Command("brew", "services", "start", "postgresql")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// startDocker starts PostgreSQL in Docker.
func (h *Helper) startDocker(port int, password string) error {
	// Check if already running
	if id := h.getDockerContainer(); id != "" {
		return fmt.Errorf("PostgreSQL container already running: %s", id)
	}

	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".postgres-data")

	// Create data directory
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	args := []string{
		"run", "-d",
		"--name", "postgres",
		"-p", fmt.Sprintf("%d:5432", port),
		"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		"-v", fmt.Sprintf("%s:/var/lib/postgresql/data", dataDir),
		"postgres:latest",
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting PostgreSQL on port %d...\n", port)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("PostgreSQL started.\n")
	fmt.Printf("Connect with: psql -h localhost -p %d -U postgres\n", port)
	return nil
}

// Stop stops PostgreSQL server.
func (h *Helper) Stop() error {
	if h.dryRun {
		fmt.Println("[dry-run] would stop PostgreSQL")
		return nil
	}

	// Check for Docker container first
	containerID := h.getDockerContainer()
	if containerID != "" {
		cmd := exec.Command("docker", "stop", containerID)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}
		exec.Command("docker", "rm", containerID).Run()
		return nil
	}

	// Stop Homebrew service
	cmd := exec.Command("brew", "services", "stop", "postgresql")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListDatabases lists all databases.
func (h *Helper) ListDatabases(host, port, user string) ([]Database, error) {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	args := []string{
		"-h", host, "-p", port, "-U", user,
		"-t", "-A", "-c",
		"SELECT datname, pg_catalog.pg_get_userbyid(datdba), pg_size_pretty(pg_database_size(datname)) FROM pg_database WHERE datistemplate = false;",
	}

	cmd := exec.Command("psql", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	var databases []Database
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			databases = append(databases, Database{
				Name:  strings.TrimSpace(parts[0]),
				Owner: strings.TrimSpace(parts[1]),
				Size:  strings.TrimSpace(parts[2]),
			})
		}
	}

	return databases, nil
}

// CreateDatabase creates a new database.
func (h *Helper) CreateDatabase(name, host, port, user string) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would create database: %s\n", name)
		return nil
	}

	args := []string{"-h", host, "-p", port, "-U", user, name}
	cmd := exec.Command("createdb", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DropDatabase drops a database.
func (h *Helper) DropDatabase(name, host, port, user string) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would drop database: %s\n", name)
		return nil
	}

	args := []string{"-h", host, "-p", port, "-U", user, name}
	cmd := exec.Command("dropdb", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Connect connects to a database.
func (h *Helper) Connect(database, host, port, user string) error {
	if database == "" {
		database = "postgres"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	args := []string{"-h", host, "-p", port, "-U", user, database}
	cmd := exec.Command("psql", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Dump dumps a database.
func (h *Helper) Dump(database, outputPath, host, port, user string) error {
	if database == "" {
		return fmt.Errorf("database name is required")
	}
	if outputPath == "" {
		outputPath = database + ".sql"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would dump %s to %s\n", database, outputPath)
		return nil
	}

	args := []string{"-h", host, "-p", port, "-U", user, "-f", outputPath, database}
	cmd := exec.Command("pg_dump", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Restore restores a database from dump.
func (h *Helper) Restore(database, inputPath, host, port, user string) error {
	if inputPath == "" {
		return fmt.Errorf("input file is required")
	}
	if database == "" {
		return fmt.Errorf("database name is required")
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would restore %s from %s\n", database, inputPath)
		return nil
	}

	args := []string{"-h", host, "-p", port, "-U", user, "-d", database, "-f", inputPath}
	cmd := exec.Command("psql", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs PostgreSQL.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("psql"); err == nil {
		return fmt.Errorf("PostgreSQL is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install PostgreSQL via Homebrew")
		return nil
	}

	fmt.Println("Installing PostgreSQL...")
	cmd := exec.Command("brew", "install", "postgresql@16")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
