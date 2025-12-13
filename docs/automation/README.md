# Automation Framework

A comprehensive automation toolkit for modern development workflows, infrastructure management, and system administration.

## 🚀 Quick Start

```bash
# Setup the automation framework
cd .automation && ./setup.sh

# Start using automation
auto --help
auto dev init python my-fastapi-app
auto k8s cluster info
auto github repo create my-new-project
```

## 🎯 Features

### 🛠️ Development Workflow Automation
- **Multi-language project initialization** (Python, Go, TypeScript)
- **FastAPI/Express/Cobra CLI templates** with best practices
- **Automated testing and building** across different project types
- **Code formatting and linting** integration
- **VS Code workspace setup** with language-specific configurations

### ☸️ Kubernetes Operations
- **Cluster management** (info, switching, monitoring)
- **Application deployment** with manifest generation
- **Helm chart management** (install, upgrade, rollback)
- **Pod operations** (logs, exec, port-forwarding)
- **Monitoring stack setup** (Prometheus + Grafana)
- **Resource cleanup** and backup operations

### 🔗 GitHub Integration  
- **Repository management** (create, clone, fork)
- **Pull request workflows** (create, merge, review)
- **Issue tracking** automation
- **CI/CD pipeline setup** (GitHub Actions templates)
- **Security scanning** and policy setup
- **Release management** automation

### 🖥️ System Administration
- **Environment setup** and configuration
- **Backup and restore** operations
- **System monitoring** and health checks
- **Security scanning** and hardening
- **Package management** automation
- **Service management** across platforms

### ⚙️ Configuration Management
- **Profile management** (development, production, etc.)
- **Template system** for quick environment setup
- **Environment variable management**
- **Configuration validation** and backup
- **Cross-machine synchronization**

### 🔐 Secrets Management
- **API key management** for all cloud providers and services
- **Interactive setup wizard** with guided configuration
- **Credential validation** and authentication testing
- **Encrypted secrets storage** with OpenSSL
- **Environment loading** and shell integration
- **Security scanning** and best practices enforcement

### ☁️ Multi-Cloud Management
- **AWS automation** (EC2, S3, Lambda, RDS, ECS, CloudFormation)
- **Azure automation** (VMs, Storage, Functions, AKS, SQL, App Service)
- **DigitalOcean automation** (Droplets, Kubernetes, Volumes, Spaces, Databases)
- **Unified cloud interface** with provider switching
- **Multi-cloud deployments** and resource comparison
- **Cost analysis** and optimization across providers

## 📋 Command Reference

### Development (`auto dev`)
```bash
auto dev init python my-api              # Create FastAPI project
auto dev init go my-cli --cobra          # Create Go CLI with Cobra
auto dev init typescript my-app          # Create TypeScript project
auto dev test                            # Run tests
auto dev build                           # Build project
auto dev format                          # Format code
```

### Kubernetes (`auto k8s`)
```bash
auto k8s cluster info                    # Show cluster information
auto k8s deploy my-app production        # Deploy application
auto k8s scale my-deployment 3           # Scale deployment
auto k8s logs my-pod --follow            # Follow pod logs
auto k8s port-forward my-pod 8080:80     # Port forward
auto k8s monitoring                      # Setup monitoring stack
auto k8s backup my-namespace             # Backup namespace
```

### GitHub (`auto github`)
```bash
auto github repo create my-project       # Create repository
auto github pr create "New feature"      # Create pull request
auto github issue create "Bug report"    # Create issue
auto github workflow setup-ci python     # Setup CI workflow
auto github security                     # Configure security
auto github release create v1.0.0        # Create release
```

### System (`auto system`)
```bash
auto system setup                        # Initial system setup
auto system backup /important/data       # Backup directory
auto system monitor                      # System monitoring
auto system cleanup                      # Clean temporary files
auto system security --scan              # Security scan
auto system packages --update            # Update packages
```

### Configuration (`auto config`)
```bash
auto config profile create work          # Create work profile
auto config profile switch development   # Switch to dev profile
auto config template apply python-dev    # Apply Python template
auto config environment create staging   # Create staging environment
auto config backup                       # Backup configuration
```

### Secrets (`auto secrets`)
```bash
auto secrets init                        # Initialize secrets management
auto secrets setup                       # Interactive setup wizard
auto secrets check-requirements          # Check missing API keys
auto secrets validate                    # Validate all credentials
auto secrets aws                         # Setup AWS credentials
auto secrets azure                       # Setup Azure credentials
auto secrets digitalocean               # Setup DigitalOcean credentials
auto secrets github                     # Setup GitHub credentials
auto secrets encrypt                    # Encrypt secrets vault
auto secrets list                       # List configured keys
```

### Cloud Management (`auto cloud`)
```bash
auto cloud status                        # Show all cloud provider status
auto cloud switch aws                    # Switch to AWS as active provider
auto cloud compare vm-sizes              # Compare instance types
auto cloud multi-deploy web-app          # Deploy to multiple clouds
auto cloud cost-compare                  # Compare pricing across providers
auto cloud backup-all                    # Backup all cloud configurations
```

### AWS (`auto aws`)
```bash
auto aws auth setup                      # Setup AWS authentication
auto aws ec2 create web-server           # Create EC2 instance
auto aws s3 create my-bucket             # Create S3 bucket
auto aws lambda deploy my-function       # Deploy Lambda function
auto aws rds create my-db mysql          # Create RDS database
auto aws cloudformation deploy my-stack  # Deploy CloudFormation stack
auto aws cost analyze                    # Analyze AWS costs
```

### Azure (`auto azure`)
```bash
auto azure auth setup                    # Setup Azure authentication
auto azure vm create web-server          # Create virtual machine
auto azure storage create mystorageacct  # Create storage account
auto azure functions create my-func-app  # Create function app
auto azure aks create my-cluster         # Create AKS cluster
auto azure sql create my-server          # Create SQL server
```

### DigitalOcean (`auto digitalocean`)
```bash
auto digitalocean auth setup             # Setup DO authentication
auto digitalocean droplets create web-1  # Create droplet
auto digitalocean kubernetes create k8s  # Create Kubernetes cluster
auto digitalocean volumes create data-vol # Create block storage volume
auto digitalocean databases create my-db # Create managed database
auto digitalocean cost analyze           # Analyze usage and costs
```

## 🏗️ Architecture

```
.automation/
├── auto                    # Main CLI entry point
├── framework/              # Core framework
│   ├── core.sh            # Logging, utilities, error handling
│   └── utils.sh           # Extended utilities
├── modules/                # Feature modules
│   ├── dev.sh             # Development workflows
│   ├── k8s.sh             # Kubernetes operations
│   ├── github.sh          # GitHub integration
│   ├── system.sh          # System administration
│   ├── config.sh          # Configuration management
│   ├── cloud.sh           # Unified cloud management
│   ├── aws.sh             # Amazon Web Services
│   ├── azure.sh           # Microsoft Azure
│   └── digitalocean.sh    # DigitalOcean
├── config/                 # Configuration files
├── logs/                   # Activity logs
├── cache/                  # Temporary files and caches
├── templates/              # Project and config templates
└── cloud/                  # Cloud-specific resources
    ├── templates/          # Multi-cloud deployment templates
    ├── configs/            # Cloud provider configurations
    └── profiles/           # Cloud environment profiles
```

## 🔧 Configuration

The framework uses configuration files in `.automation/config/`:

- **`automation.conf`** - Main framework settings
- **`profiles/`** - Environment profiles (work, development, etc.)
- **`templates/`** - Project and configuration templates
- **`environments/`** - Environment-specific variables

### Example Configuration

```bash
# automation.conf
AUTO_LOG_LEVEL=INFO
DEV_PROJECTS_DIR=$HOME/projects
GITHUB_DEFAULT_VISIBILITY=public
K8S_DEFAULT_NAMESPACE=default
```

## 🚀 Advanced Usage

### Creating Custom Templates

```bash
# Create a custom project template
auto config template create my-template custom

# Apply template to current environment
auto config template apply my-template
```

### Profile Management

```bash
# Create profiles for different environments
auto config profile create work "Work environment"
auto config profile create personal "Personal projects"

# Switch between profiles
auto config profile switch work
```

### Batch Operations

```bash
# Deploy multiple applications
for app in api worker scheduler; do
    auto k8s deploy $app production
done

# Create multiple repositories
for repo in frontend backend mobile; do
    auto github repo create $repo
done
```

## 🔌 Integration

### Shell Integration

The framework integrates with your shell environment through `.bash_tools/automation.sh`:

```bash
# Quick aliases available in your shell
ainit python my-api     # Quick project init
adeploy my-app prod     # Quick deployment
apr "Feature X"         # Quick PR creation
abackup ~/important     # Quick backup
```

### VS Code Integration

Automatically configures VS Code workspaces with:
- Language-specific settings
- Debugging configurations
- Recommended extensions
- Task runners

### CI/CD Integration

Generates GitHub Actions workflows for:
- Multi-platform testing
- Security scanning
- Automated releases
- Container building

## 🔒 Security

- **Secrets management** through encrypted environment files
- **Security scanning** integration with CodeQL and vulnerability alerts
- **Permission validation** for system operations
- **Audit logging** of all automation activities

## 🐛 Troubleshooting

### Common Issues

1. **Permission denied errors**
   ```bash
   chmod +x .automation/auto
   .automation/setup.sh
   ```

2. **Command not found**
   ```bash
   export PATH=".automation:$PATH"
   source ~/.bash_profile
   ```

3. **Configuration validation**
   ```bash
   auto config validate
   ```

### Debug Mode

```bash
auto --verbose <command>  # Enable debug logging
tail -f .automation/logs/automation.log  # Follow logs
```

## 🤝 Contributing

The automation framework is designed to be extensible:

1. **Add new modules** in `.automation/modules/`
2. **Create templates** in `.automation/templates/`
3. **Extend existing modules** with new commands
4. **Add completion** for new commands

### Module Template

```bash
#!/bin/bash
# New Module Template

module_help() {
    cat << EOF
My New Module
USAGE: auto mymodule <command>
COMMANDS:
    action1    Description of action1
EOF
}

mymodule_main() {
    case "$1" in
        "help") module_help ;;
        "action1") echo "Action 1 executed" ;;
        *) module_help; exit 1 ;;
    esac
}
```

## 📈 Roadmap

- [ ] **Cloud provider integration** (AWS, GCP, Azure)
- [ ] **Container orchestration** beyond Kubernetes
- [ ] **Infrastructure as Code** (Terraform, Pulumi)
- [ ] **Monitoring and alerting** integrations
- [ ] **Database management** automation
- [ ] **ML/AI workflow** automation

---

**Built for developers, by developers.** The automation framework grows with your needs and adapts to your workflow.