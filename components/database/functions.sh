#!/bin/sh
# components/database/functions.sh - Database helper functions

# =============================================================================
# Quick Connect Helpers
# =============================================================================

# PostgreSQL local connection
pglocal() {
    pgcli -h localhost -U "${1:-postgres}" "${2:-postgres}"
}

# MySQL local connection
mylocal() {
    mycli -h localhost -u "${1:-root}" "${2:-}"
}

# MongoDB local connection
mongolocal() {
    mongosh "mongodb://localhost:27017/${1:-test}"
}

# Redis local connection
rdlocal() {
    if command -v iredis >/dev/null 2>&1; then
        iredis -h localhost -p "${1:-6379}"
    else
        redis-cli -h localhost -p "${1:-6379}"
    fi
}

# Neo4j local connection
neolocal() {
    cypher-shell -u "${1:-neo4j}" -p "${2:-neo4j}" -a "bolt://localhost:7687"
}

# =============================================================================
# Database Status Helper
# =============================================================================

db_status() {
    echo "Database Services Status"
    echo "========================"
    echo ""

    # PostgreSQL
    printf "PostgreSQL: "
    if command -v pg_isready >/dev/null 2>&1; then
        if pg_isready -q 2>/dev/null; then
            echo "Running"
        else
            echo "Not running"
        fi
    else
        echo "Not installed"
    fi

    # MySQL
    printf "MySQL:      "
    if command -v mysqladmin >/dev/null 2>&1; then
        if mysqladmin ping -u root --silent 2>/dev/null; then
            echo "Running"
        else
            echo "Not running"
        fi
    elif command -v mysql >/dev/null 2>&1; then
        echo "Installed (status unknown)"
    else
        echo "Not installed"
    fi

    # MongoDB
    printf "MongoDB:    "
    if command -v mongosh >/dev/null 2>&1; then
        if mongosh --eval "db.runCommand({ping:1})" --quiet 2>/dev/null | grep -q "ok"; then
            echo "Running"
        else
            echo "Not running"
        fi
    else
        echo "Not installed"
    fi

    # Redis
    printf "Redis:      "
    if command -v redis-cli >/dev/null 2>&1; then
        if redis-cli ping 2>/dev/null | grep -q "PONG"; then
            echo "Running"
        else
            echo "Not running"
        fi
    else
        echo "Not installed"
    fi

    # Neo4j
    printf "Neo4j:      "
    if command -v neo4j >/dev/null 2>&1; then
        if neo4j status 2>/dev/null | grep -q "running"; then
            echo "Running"
        else
            echo "Not running"
        fi
    elif command -v cypher-shell >/dev/null 2>&1; then
        echo "Installed (status unknown)"
    else
        echo "Not installed"
    fi

    echo ""
}

# =============================================================================
# Service Management (macOS)
# =============================================================================

# Start all common databases
db_start_all() {
    if [ "$CURRENT_PLATFORM" != "darwin" ]; then
        echo "This function is for macOS with Homebrew"
        return 1
    fi

    echo "Starting database services..."
    brew services start postgresql@14 2>/dev/null
    brew services start redis 2>/dev/null
    brew services start mongodb-community 2>/dev/null
    echo "Done. Use 'db_status' to check status."
}

# Stop all common databases
db_stop_all() {
    if [ "$CURRENT_PLATFORM" != "darwin" ]; then
        echo "This function is for macOS with Homebrew"
        return 1
    fi

    echo "Stopping database services..."
    brew services stop postgresql@14 2>/dev/null
    brew services stop redis 2>/dev/null
    brew services stop mongodb-community 2>/dev/null
    echo "Done."
}
