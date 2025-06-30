#!/bin/bash
# Test Configuration and Utilities
# Common test functions and configuration for the test suite

# Test configuration
export TEST_TIMEOUT=30
export TEST_VERBOSE=${TEST_VERBOSE:-false}
export TEST_SKIP_CLOUD=${TEST_SKIP_CLOUD:-false}
export TEST_SKIP_INTEGRATION=${TEST_SKIP_INTEGRATION:-false}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test utilities
log_test() {
    local level="$1"
    shift
    local message="$*"
    
    case "$level" in
        "INFO")  echo -e "${BLUE}[TEST-INFO]${NC} $message" ;;
        "PASS")  echo -e "${GREEN}[TEST-PASS]${NC} $message" ;;
        "FAIL")  echo -e "${RED}[TEST-FAIL]${NC} $message" ;;
        "WARN")  echo -e "${YELLOW}[TEST-WARN]${NC} $message" ;;
        "SKIP")  echo -e "${YELLOW}[TEST-SKIP]${NC} $message" ;;
    esac
}

# Test assertion functions
assert_command_exists() {
    local cmd="$1"
    if command -v "$cmd" >/dev/null 2>&1; then
        log_test "PASS" "Command '$cmd' exists"
        return 0
    else
        log_test "FAIL" "Command '$cmd' not found"
        return 1
    fi
}

assert_file_exists() {
    local file="$1"
    if [ -f "$file" ]; then
        log_test "PASS" "File '$file' exists"
        return 0
    else
        log_test "FAIL" "File '$file' not found"
        return 1
    fi
}

assert_directory_exists() {
    local dir="$1"
    if [ -d "$dir" ]; then
        log_test "PASS" "Directory '$dir' exists"
        return 0
    else
        log_test "FAIL" "Directory '$dir' not found"
        return 1
    fi
}

assert_function_defined() {
    local func="$1"
    if declare -f "$func" >/dev/null 2>&1; then
        log_test "PASS" "Function '$func' is defined"
        return 0
    else
        log_test "FAIL" "Function '$func' not defined"
        return 1
    fi
}

assert_variable_set() {
    local var="$1"
    if [ -n "${!var}" ]; then
        log_test "PASS" "Variable '$var' is set"
        return 0
    else
        log_test "FAIL" "Variable '$var' not set"
        return 1
    fi
}

assert_alias_exists() {
    local alias_name="$1"
    if alias "$alias_name" >/dev/null 2>&1; then
        log_test "PASS" "Alias '$alias_name' exists"
        return 0
    else
        log_test "FAIL" "Alias '$alias_name' not found"
        return 1
    fi
}

# Test execution wrapper
run_test() {
    local test_name="$1"
    local test_function="$2"
    
    log_test "INFO" "Running test: $test_name"
    
    if [ "$TEST_VERBOSE" = "true" ]; then
        if $test_function; then
            log_test "PASS" "Test '$test_name' passed"
            return 0
        else
            log_test "FAIL" "Test '$test_name' failed"
            return 1
        fi
    else
        if $test_function >/dev/null 2>&1; then
            log_test "PASS" "Test '$test_name' passed"
            return 0
        else
            log_test "FAIL" "Test '$test_name' failed"
            return 1
        fi
    fi
}

# Test environment detection
detect_test_environment() {
    log_test "INFO" "Detecting test environment..."
    
    # OS Detection
    case "$(uname -s)" in
        Darwin) export TEST_OS="macos" ;;
        Linux)  export TEST_OS="linux" ;;
        *)      export TEST_OS="unknown" ;;
    esac
    log_test "INFO" "Operating System: $TEST_OS"
    
    # Shell Detection
    if [ -n "$BASH_VERSION" ]; then
        export TEST_SHELL="bash"
    elif [ -n "$ZSH_VERSION" ]; then
        export TEST_SHELL="zsh"
    else
        export TEST_SHELL="unknown"
    fi
    log_test "INFO" "Shell: $TEST_SHELL"
    
    # CI Detection
    if [ -n "$CI" ] || [ -n "$GITHUB_ACTIONS" ] || [ -n "$TRAVIS" ] || [ -n "$JENKINS_URL" ]; then
        export TEST_CI="true"
        log_test "INFO" "CI environment detected"
    else
        export TEST_CI="false"
        log_test "INFO" "Local environment detected"
    fi
    
    # Cloud CLI Detection
    log_test "INFO" "Checking cloud CLI availability..."
    command -v aws >/dev/null 2>&1 && export TEST_HAS_AWS="true" || export TEST_HAS_AWS="false"
    command -v az >/dev/null 2>&1 && export TEST_HAS_AZURE="true" || export TEST_HAS_AZURE="false"
    command -v doctl >/dev/null 2>&1 && export TEST_HAS_DOCTL="true" || export TEST_HAS_DOCTL="false"
    command -v kubectl >/dev/null 2>&1 && export TEST_HAS_KUBECTL="true" || export TEST_HAS_KUBECTL="false"
    command -v helm >/dev/null 2>&1 && export TEST_HAS_HELM="true" || export TEST_HAS_HELM="false"
    command -v gh >/dev/null 2>&1 && export TEST_HAS_GH="true" || export TEST_HAS_GH="false"
    
    log_test "INFO" "AWS CLI: $TEST_HAS_AWS"
    log_test "INFO" "Azure CLI: $TEST_HAS_AZURE"
    log_test "INFO" "doctl: $TEST_HAS_DOCTL"
    log_test "INFO" "kubectl: $TEST_HAS_KUBECTL"
    log_test "INFO" "helm: $TEST_HAS_HELM"
    log_test "INFO" "gh: $TEST_HAS_GH"
}

# Test data generation
create_test_project() {
    local project_name="$1"
    local project_type="${2:-python}"
    local test_dir="tests/test-projects/$project_name"
    
    mkdir -p "$test_dir"
    
    case "$project_type" in
        "python")
            cat > "$test_dir/main.py" << 'EOF'
def hello():
    return "Hello, World!"

if __name__ == "__main__":
    print(hello())
EOF
            cat > "$test_dir/requirements.txt" << 'EOF'
# Test requirements
requests>=2.25.0
EOF
            ;;
        "go")
            cat > "$test_dir/main.go" << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF
            cat > "$test_dir/go.mod" << EOF
module $project_name

go 1.19
EOF
            ;;
        "typescript")
            cat > "$test_dir/package.json" << EOF
{
    "name": "$project_name",
    "version": "1.0.0",
    "main": "index.js",
    "scripts": {
        "test": "echo \"Error: no test specified\" && exit 1"
    }
}
EOF
            cat > "$test_dir/index.ts" << 'EOF'
function hello(): string {
    return "Hello, World!";
}

console.log(hello());
EOF
            ;;
    esac
    
    log_test "INFO" "Created test project: $test_dir ($project_type)"
}

# Cleanup functions
cleanup_test_environment() {
    log_test "INFO" "Cleaning up test environment..."
    
    # Remove test projects
    rm -rf tests/test-projects
    
    # Clean up any test virtual environments
    if [ -d "tests/test-venvs" ]; then
        rm -rf tests/test-venvs
    fi
    
    # Clean up test files
    find tests -name "*.tmp" -delete 2>/dev/null || true
    find tests -name "*.test" -delete 2>/dev/null || true
    
    log_test "INFO" "Test environment cleaned up"
}

# Performance measurement
measure_time() {
    local description="$1"
    shift
    local command="$@"
    
    log_test "INFO" "Measuring: $description"
    
    local start_time=$(date +%s.%N)
    eval "$command"
    local end_time=$(date +%s.%N)
    
    local duration=$(echo "$end_time - $start_time" | bc 2>/dev/null || echo "unknown")
    log_test "INFO" "Duration: ${duration}s"
}

# Initialize test environment
init_test_environment() {
    detect_test_environment
    
    # Create test directories
    mkdir -p tests/{logs,backups,test-projects,temp}
    
    # Set test-specific environment variables
    export AUTO_DRY_RUN=true
    export AUTO_LOG_LEVEL=DEBUG
    
    log_test "INFO" "Test environment initialized"
}

# Test suite runner
run_test_suite() {
    local suite_name="$1"
    shift
    local tests=("$@")
    
    log_test "INFO" "Running test suite: $suite_name"
    
    local passed=0
    local failed=0
    local skipped=0
    
    for test in "${tests[@]}"; do
        if run_test "$test" "$test"; then
            ((passed++))
        else
            if [[ "$test" =~ skip ]]; then
                ((skipped++))
            else
                ((failed++))
            fi
        fi
    done
    
    log_test "INFO" "Test suite '$suite_name' completed"
    log_test "INFO" "Results: $passed passed, $failed failed, $skipped skipped"
    
    return $failed
}

# Export functions for use in other test scripts
export -f log_test assert_command_exists assert_file_exists assert_directory_exists
export -f assert_function_defined assert_variable_set assert_alias_exists
export -f run_test detect_test_environment create_test_project cleanup_test_environment
export -f measure_time init_test_environment run_test_suite