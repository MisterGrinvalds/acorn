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

// Deployment represents a kubernetes deployment.
type Deployment struct {
	Name      string `json:"name" yaml:"name"`
	Ready     string `json:"ready" yaml:"ready"`
	UpToDate  string `json:"up_to_date" yaml:"up_to_date"`
	Available string `json:"available" yaml:"available"`
	Age       string `json:"age" yaml:"age"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// Service represents a kubernetes service.
type Service struct {
	Name       string `json:"name" yaml:"name"`
	Type       string `json:"type" yaml:"type"`
	ClusterIP  string `json:"cluster_ip" yaml:"cluster_ip"`
	ExternalIP string `json:"external_ip,omitempty" yaml:"external_ip,omitempty"`
	Ports      string `json:"ports" yaml:"ports"`
	Age        string `json:"age" yaml:"age"`
	Namespace  string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// Event represents a kubernetes event.
type Event struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	LastSeen  string `json:"last_seen" yaml:"last_seen"`
	Type      string `json:"type" yaml:"type"`
	Reason    string `json:"reason" yaml:"reason"`
	Object    string `json:"object" yaml:"object"`
	Message   string `json:"message" yaml:"message"`
}

// GetDeployments returns list of deployments.
func (h *Helper) GetDeployments(namespace string) ([]Deployment, error) {
	args := []string{"get", "deployments", "-o", "jsonpath={range .items[*]}{.metadata.name},{.status.readyReplicas}/{.spec.replicas},{.status.updatedReplicas},{.status.availableReplicas},{.metadata.creationTimestamp}{\"\\n\"}{end}"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	var deployments []Deployment
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 5)
		d := Deployment{Name: parts[0], Namespace: namespace}
		if len(parts) > 1 {
			d.Ready = parts[1]
		}
		if len(parts) > 2 {
			d.UpToDate = parts[2]
		}
		if len(parts) > 3 {
			d.Available = parts[3]
		}
		if len(parts) > 4 {
			d.Age = parts[4]
		}
		deployments = append(deployments, d)
	}

	return deployments, nil
}

// GetServices returns list of services.
func (h *Helper) GetServices(namespace string) ([]Service, error) {
	args := []string{"get", "services", "-o", "jsonpath={range .items[*]}{.metadata.name},{.spec.type},{.spec.clusterIP},{.spec.externalIPs[0]},{.spec.ports[*].port},{.metadata.creationTimestamp}{\"\\n\"}{end}"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	var services []Service
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 6)
		s := Service{Name: parts[0], Namespace: namespace}
		if len(parts) > 1 {
			s.Type = parts[1]
		}
		if len(parts) > 2 {
			s.ClusterIP = parts[2]
		}
		if len(parts) > 3 {
			s.ExternalIP = parts[3]
		}
		if len(parts) > 4 {
			s.Ports = parts[4]
		}
		if len(parts) > 5 {
			s.Age = parts[5]
		}
		services = append(services, s)
	}

	return services, nil
}

// GetEvents returns recent events sorted by time.
func (h *Helper) GetEvents(namespace string) ([]Event, error) {
	args := []string{"get", "events", "--sort-by=.lastTimestamp", "-o", "jsonpath={range .items[*]}{.metadata.namespace},{.lastTimestamp},{.type},{.reason},{.involvedObject.kind}/{.involvedObject.name},{.message}{\"\\n\"}{end}"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var events []Event
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 6)
		e := Event{}
		if len(parts) > 0 {
			e.Namespace = parts[0]
		}
		if len(parts) > 1 {
			e.LastSeen = parts[1]
		}
		if len(parts) > 2 {
			e.Type = parts[2]
		}
		if len(parts) > 3 {
			e.Reason = parts[3]
		}
		if len(parts) > 4 {
			e.Object = parts[4]
		}
		if len(parts) > 5 {
			e.Message = parts[5]
		}
		events = append(events, e)
	}

	return events, nil
}

// RolloutStatus gets the rollout status of a deployment.
func (h *Helper) RolloutStatus(deployment, namespace string) error {
	args := []string{"rollout", "status", "deployment/" + deployment}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RolloutRestart performs a rolling restart of a deployment.
func (h *Helper) RolloutRestart(deployment, namespace string) error {
	args := []string{"rollout", "restart", "deployment/" + deployment}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RolloutUndo undoes a rollout to a previous revision.
func (h *Helper) RolloutUndo(deployment, namespace string, revision int) error {
	args := []string{"rollout", "undo", "deployment/" + deployment}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if revision > 0 {
		args = append(args, "--to-revision", fmt.Sprintf("%d", revision))
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ScaleDeployment scales a deployment to the specified replicas.
func (h *Helper) ScaleDeployment(deployment, namespace string, replicas int) error {
	args := []string{"scale", "deployment/" + deployment, fmt.Sprintf("--replicas=%d", replicas)}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetLogs gets logs from a pod.
func (h *Helper) GetLogs(pod, namespace string, follow bool, tail int) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, pod)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecInPod executes a command in a pod.
func (h *Helper) ExecInPod(pod, namespace string, command []string) error {
	args := []string{"exec", "-it"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, pod, "--")
	args = append(args, command...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PortForward forwards a local port to a pod port.
func (h *Helper) PortForward(pod, namespace, ports string) error {
	args := []string{"port-forward"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, pod, ports)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DescribePod describes a pod.
func (h *Helper) DescribePod(pod, namespace string) error {
	args := []string{"describe", "pod"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, pod)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DeletePod deletes a pod.
func (h *Helper) DeletePod(pod, namespace string) error {
	args := []string{"delete", "pod"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, pod)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ApplyFile applies a manifest file.
func (h *Helper) ApplyFile(filepath string, dryRunMode string) error {
	args := []string{"apply", "-f", filepath}
	if dryRunMode != "" {
		args = append(args, "--dry-run="+dryRunMode)
	}

	if h.dryRun && dryRunMode == "" {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DeleteFile deletes resources from a manifest file.
func (h *Helper) DeleteFile(filepath string) error {
	args := []string{"delete", "-f", filepath}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TopPods returns resource usage for pods.
func (h *Helper) TopPods(namespace string) error {
	args := []string{"top", "pods"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TopNodes returns resource usage for nodes.
func (h *Helper) TopNodes() error {
	cmd := exec.Command("kubectl", "top", "nodes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// WatchResource watches a resource type.
func (h *Helper) WatchResource(resource, namespace string) error {
	args := []string{"get", resource, "-w"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetAsYAML returns a resource as YAML.
func (h *Helper) GetAsYAML(resource, name, namespace string) (string, error) {
	args := []string{"get", resource}
	if name != "" {
		args = append(args, name)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "-o", "yaml")

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", resource, err)
	}

	return string(out), nil
}

// GetAsJSON returns a resource as JSON.
func (h *Helper) GetAsJSON(resource, name, namespace string) (string, error) {
	args := []string{"get", resource}
	if name != "" {
		args = append(args, name)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	args = append(args, "-o", "json")

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", resource, err)
	}

	return string(out), nil
}
