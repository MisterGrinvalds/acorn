# 🔐 Secrets Management

Comprehensive secrets and API key management for the automation framework. Securely store, validate, and manage API keys for cloud providers, development tools, and services.

## 🚀 Quick Start

```bash
# Initialize secrets management
auto secrets init

# Interactive setup wizard
auto secrets setup

# Check what's configured
auto secrets check-requirements

# Validate all API keys
auto secrets validate

# Load secrets into current shell
load_secrets
```

## 📋 Supported Services

### ☁️ Cloud Providers
- **AWS** - Access Key ID, Secret Access Key, Region, Session Token
- **Azure** - Client ID, Client Secret, Tenant ID, Subscription ID  
- **DigitalOcean** - API Token, Spaces Access Key

### 🔗 Development Tools
- **GitHub** - Personal Access Token, Username, Organization
- **GitLab** - API Token, Username
- **Docker Hub** - Username, Password, Email

### ☸️ Container & Orchestration
- **GitHub Container Registry** - Token
- **Azure Container Registry** - Username, Password
- **Kubernetes** - Service Account Token, CA Certificate, Cluster Endpoint
- **Helm** - Repository Username, Password

### 🗄️ Databases
- **PostgreSQL** - Username, Password, Host, Port, Database
- **MySQL** - Username, Password, Host, Port, Database
- **MongoDB** - URI, Username, Password
- **Redis** - URL, Password

### 📊 Monitoring & Observability
- **DataDog** - API Key, App Key
- **New Relic** - License Key, API Key
- **Grafana** - API Token
- **Prometheus** - Username, Password

### 💬 Communication
- **Slack** - Bot Token, App Token, Webhook URL
- **Discord** - Bot Token, Webhook URL

### 🔄 CI/CD
- **Jenkins** - URL, Username, API Token
- **CircleCI** - Token
- **Travis CI** - Token

### 🛠️ Development Services
- **Sentry** - DSN, Auth Token
- **OpenAI** - API Key
- **Anthropic** - API Key

## 🔧 Setup Methods

### Interactive Setup Wizard

```bash
# Full setup wizard (recommended for first-time setup)
auto secrets setup

# Individual provider setup
auto secrets aws
auto secrets azure
auto secrets digitalocean
auto secrets github
```

### Manual Configuration

```bash
# Initialize secrets management
auto secrets init

# Edit the secrets file directly
nano ~/.automation/secrets/.env

# Validate your changes
auto secrets validate
```

### Environment Variables

Set environment variables directly:

```bash
export AWS_ACCESS_KEY_ID="your-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export GITHUB_TOKEN="your-github-token"
```

## 🔍 Validation & Testing

### Check Configuration Status

```bash
# Check which API keys are configured
auto secrets check-requirements

# List configured keys (without values)
auto secrets list

# Check specific provider
check-aws        # Check if AWS credentials are loaded
check-github     # Check if GitHub token is loaded
check-azure      # Check if Azure credentials are loaded
check-do         # Check if DigitalOcean token is loaded
```

### Validate API Keys

```bash
# Validate all configured credentials
auto secrets validate

# Quick validation from shell
validate_secrets
```

### Testing Integration

Use the Makefile testing framework:

```bash
# Test API key configuration
make test-api-keys

# Test authentication status for all services
make test-auth-status

# Check required CLI tools
make test-required-tools

# Run secrets setup wizard
make setup-secrets-wizard
```

## 🔄 Daily Usage

### Loading Secrets

```bash
# Load secrets into current shell environment
load_secrets

# Auto-load on shell startup (add to .bash_profile)
export AUTO_LOAD_SECRETS=true
```

### Quick Status Checks

```bash
# Quick status check
secrets-status

# Validate credentials
secrets-validate

# List what's configured
secrets-list
```

### Shell Integration

Available aliases and functions:
- `load_secrets` / `secrets-load` - Load secrets into environment
- `secrets_status` / `secrets-status` - Check configuration status
- `validate_secrets` / `secrets-validate` - Validate all credentials
- `setup-aws`, `setup-azure`, `setup-github`, `setup-do` - Quick setup
- `check-aws`, `check-azure`, `check-github`, `check-do` - Check specific keys

## 🔒 Security Features

### Encryption

```bash
# Encrypt secrets vault
auto secrets encrypt

# Decrypt secrets vault
auto secrets decrypt
```

### File Permissions

Secrets files are automatically secured:
- `.automation/secrets/` directory: `700` (owner only)
- Secrets files: `600` (owner read/write only)
- Git exclusion: Automatic `.gitignore` creation

### Best Practices

1. **Never commit secrets to version control**
2. **Use encrypted vault for sensitive environments**
3. **Regularly rotate API keys**
4. **Validate credentials periodically**
5. **Use least-privilege access**

## 📁 File Structure

```
.automation/secrets/
├── .gitignore          # Protects against accidental commits
├── config              # Secrets management configuration
├── template.env        # Template for new secrets
├── .env                # Your actual secrets (never commit!)
└── vault.enc           # Encrypted secrets vault (optional)
```

## 🔑 API Key Requirements by Provider

### AWS
- **Required**: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- **Optional**: `AWS_DEFAULT_REGION`, `AWS_SESSION_TOKEN`, `AWS_PROFILE`
- **CLI**: `aws configure` or environment variables

### Azure  
- **Required**: `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`
- **Optional**: `AZURE_SUBSCRIPTION_ID`
- **CLI**: `az login --service-principal`

### DigitalOcean
- **Required**: `DIGITALOCEAN_TOKEN`
- **Optional**: `DIGITALOCEAN_SPACES_ACCESS_KEY`, `DIGITALOCEAN_SPACES_SECRET_KEY`
- **CLI**: `doctl auth init`

### GitHub
- **Required**: `GITHUB_TOKEN`
- **Optional**: `GITHUB_USERNAME`, `GITHUB_ORGANIZATION`
- **CLI**: `gh auth login`

## 🧪 Testing Your Setup

### Makefile Targets

```bash
# Check API key configuration
make test-api-keys-check

# Validate configured API keys  
make test-api-keys-validate

# Test authentication status
make test-auth-status

# Check required CLI tools
make test-required-tools

# Interactive secrets setup
make setup-secrets-wizard
```

### Manual Testing

```bash
# Test AWS
aws sts get-caller-identity

# Test Azure
az account show

# Test DigitalOcean
doctl account get

# Test GitHub
gh auth status

# Test Kubernetes
kubectl cluster-info
```

## 🔄 Common Workflows

### Initial Setup

```bash
# 1. Initialize secrets management
auto secrets init

# 2. Run interactive setup
auto secrets setup

# 3. Load secrets into environment
load_secrets

# 4. Validate everything works
auto secrets validate
```

### Adding New Provider

```bash
# 1. Edit template if needed
auto secrets template > new-provider.env

# 2. Add your credentials
nano ~/.automation/secrets/.env

# 3. Validate new credentials
auto secrets validate
```

### Rotating Keys

```bash
# 1. Update credentials in provider console
# 2. Update local secrets file
auto secrets setup  # Re-run setup for specific provider

# 3. Validate new credentials
auto secrets validate

# 4. Test with actual API calls
make test-auth-status
```

### Troubleshooting

```bash
# Check what's missing
auto secrets check-requirements

# Validate specific provider
check-aws && echo "AWS OK" || echo "AWS needs setup"

# View detailed logs
tail -f ~/.automation/logs/automation.log

# Reset and start over
rm ~/.automation/secrets/.env
auto secrets setup
```

## 🔗 Integration Examples

### Use in Scripts

```bash
#!/bin/bash
# Load secrets at start of script
source ~/.automation/secrets/.env

# Or use the function
load_secrets

# Now use the credentials
aws s3 ls
az vm list
doctl droplets list
```

### Use with Automation

```bash
# Deploy with credentials loaded
load_secrets
auto cloud deploy my-app production

# Backup with proper authentication
auto system backup /data --cloud aws
```

### CI/CD Integration

```bash
# GitHub Actions
env:
  AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
  AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}

# Use in workflow
- name: Test credentials
  run: make test-auth-status
```

## 📚 Advanced Features

### Profiles

```bash
# Create environment-specific profiles
auto config profile create production
auto config profile create development

# Switch between profiles
auto config profile switch production
```

### Templates

```bash
# Create custom templates
auto config template create my-stack

# Apply templates
auto config template apply my-stack
```

### Backup & Restore

```bash
# Backup secrets
auto secrets backup

# Restore from backup
auto secrets restore backup-20240101.tar.gz
```

---

**Security Note**: Always treat API keys as sensitive information. Never share them, commit them to version control, or store them in plain text in unsecured locations. Use the encryption features for additional security in shared environments.