// Package kubernetes provides Kubernetes helper functionality.
package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ContextInfo represents current kubernetes context info.
type ContextInfo struct {
	Context   string `json:"context" yaml:"context"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Server    string `json:"server" yaml:"server"`
}

// Context represents a kubernetes context.
type Context struct {
	Name    string `json:"name" yaml:"name"`
	Cluster string `json:"cluster" yaml:"cluster"`
	User    string `json:"user" yaml:"user"`
	Current bool   `json:"current" yaml:"current"`
}

// Namespace represents a kubernetes namespace.
type Namespace struct {
	Name   string `json:"name" yaml:"name"`
	Status string `json:"status" yaml:"status"`
}

// Pod represents a kubernetes pod.
type Pod struct {
	Name      string `json:"name" yaml:"name"`
	Ready     string `json:"ready" yaml:"ready"`
	Status    string `json:"status" yaml:"status"`
	Restarts  string `json:"restarts" yaml:"restarts"`
	Age       string `json:"age" yaml:"age"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// Helper provides Kubernetes helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Kubernetes Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsKubectlInstalled checks if kubectl is installed.
func (h *Helper) IsKubectlInstalled() bool {
	_, err := exec.LookPath("kubectl")
	return err == nil
}

// GetContextInfo returns current context info.
func (h *Helper) GetContextInfo() (*ContextInfo, error) {
	info := &ContextInfo{}

	// Get current context
	cmd := exec.Command("kubectl", "config", "current-context")
	if out, err := cmd.Output(); err == nil {
		info.Context = strings.TrimSpace(string(out))
	}

	// Get namespace
	cmd = exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={..namespace}")
	if out, err := cmd.Output(); err == nil {
		info.Namespace = strings.TrimSpace(string(out))
		if info.Namespace == "" {
			info.Namespace = "default"
		}
	}

	// Get server
	cmd = exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={.clusters[0].cluster.server}")
	if out, err := cmd.Output(); err == nil {
		info.Server = strings.TrimSpace(string(out))
	}

	return info, nil
}

// GetContexts returns list of contexts.
func (h *Helper) GetContexts() ([]Context, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}

	currentCmd := exec.Command("kubectl", "config", "current-context")
	currentOut, _ := currentCmd.Output()
	current := strings.TrimSpace(string(currentOut))

	var contexts []Context
	for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if name == "" {
			continue
		}
		contexts = append(contexts, Context{
			Name:    name,
			Current: name == current,
		})
	}

	return contexts, nil
}

// UseContext switches to a context.
func (h *Helper) UseContext(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl config use-context %s\n", name)
		return nil
	}

	cmd := exec.Command("kubectl", "config", "use-context", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetNamespaces returns list of namespaces.
func (h *Helper) GetNamespaces() ([]Namespace, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "jsonpath={range .items[*]}{.metadata.name},{.status.phase}{\"\\n\"}{end}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	var namespaces []Namespace
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		ns := Namespace{Name: parts[0]}
		if len(parts) > 1 {
			ns.Status = parts[1]
		}
		namespaces = append(namespaces, ns)
	}

	return namespaces, nil
}

// UseNamespace switches to a namespace.
func (h *Helper) UseNamespace(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl config set-context --current --namespace=%s\n", name)
		return nil
	}

	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace="+name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetPods returns list of pods with optional filter.
func (h *Helper) GetPods(filter string) ([]Pod, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-o", "jsonpath={range .items[*]}{.metadata.name},{.status.containerStatuses[0].ready}/{len .status.containerStatuses},{.status.phase},{.status.containerStatuses[0].restartCount},{.metadata.creationTimestamp}{\"\\n\"}{end}")
	out, err := cmd.Output()
	if err != nil {
		// Fall back to simple output
		cmd = exec.Command("kubectl", "get", "pods")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return nil, cmd.Run()
	}

	var pods []Pod
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		if filter != "" && !strings.Contains(line, filter) {
			continue
		}
		parts := strings.SplitN(line, ",", 5)
		pod := Pod{Name: parts[0]}
		if len(parts) > 1 {
			pod.Ready = parts[1]
		}
		if len(parts) > 2 {
			pod.Status = parts[2]
		}
		if len(parts) > 3 {
			pod.Restarts = parts[3]
		}
		if len(parts) > 4 {
			pod.Age = parts[4]
		}
		pods = append(pods, pod)
	}

	return pods, nil
}

// GetAllResources returns all resources in a namespace.
func (h *Helper) GetAllResources(namespace string) error {
	if namespace == "" {
		// Get current namespace
		cmd := exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={..namespace}")
		if out, err := cmd.Output(); err == nil {
			namespace = strings.TrimSpace(string(out))
		}
		if namespace == "" {
			namespace = "default"
		}
	}

	fmt.Println("=== Pods ===")
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println()
	fmt.Println("=== Services ===")
	cmd = exec.Command("kubectl", "get", "services", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println()
	fmt.Println("=== Deployments ===")
	cmd = exec.Command("kubectl", "get", "deployments", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

// CleanEvictedPods deletes all evicted pods.
func (h *Helper) CleanEvictedPods() (int, error) {
	// Get evicted pods
	cmd := exec.Command("kubectl", "get", "pods", "--all-namespaces", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get pods: %w", err)
	}

	var result struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Status struct {
				Reason string `json:"reason"`
			} `json:"status"`
		} `json:"items"`
	}

	if err := json.Unmarshal(out, &result); err != nil {
		return 0, fmt.Errorf("failed to parse pods: %w", err)
	}

	count := 0
	for _, pod := range result.Items {
		if pod.Status.Reason == "Evicted" {
			if h.dryRun {
				fmt.Printf("[dry-run] would delete pod %s/%s\n", pod.Metadata.Namespace, pod.Metadata.Name)
				count++
				continue
			}

			delCmd := exec.Command("kubectl", "delete", "pod", "-n", pod.Metadata.Namespace, pod.Metadata.Name)
			if err := delCmd.Run(); err == nil {
				count++
			}
		}
	}

	return count, nil
}
