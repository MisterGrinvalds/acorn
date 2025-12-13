#!/bin/sh
# Database tool aliases and helper functions
# Requires: shell/discovery.sh

# =============================================================================
# PostgreSQL
# =============================================================================
alias pg='pgcli'
alias psqlc='psql'  # Direct psql client

# Brew service management (macOS)
if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias pgstart='brew services start postgresql@14'
    alias pgstop='brew services stop postgresql@14'
    alias pgrestart='brew services restart postgresql@14'
    alias pgstatus='brew services info postgresql@14'
fi

# Quick connect helpers
pglocal() {
    pgcli -h localhost -U "${1:-postgres}" "${2:-postgres}"
}

# =============================================================================
# MySQL
# =============================================================================
if command -v mycli >/dev/null 2>&1; then
    alias my='mycli'
fi

if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias mystart='brew services start mysql'
    alias mystop='brew services stop mysql'
    alias myrestart='brew services restart mysql'
    alias mystatus='brew services info mysql'
fi

# Quick connect helpers
mylocal() {
    mycli -h localhost -u "${1:-root}" "${2:-}"
}

# =============================================================================
# MongoDB
# =============================================================================
if command -v mongosh >/dev/null 2>&1; then
    alias mongo='mongosh'
    alias msh='mongosh'
fi

if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias mongostart='brew services start mongodb-community'
    alias mongostop='brew services stop mongodb-community'
    alias mongorestart='brew services restart mongodb-community'
    alias mongostatus='brew services info mongodb-community'
fi

# Quick connect helpers
mongolocal() {
    mongosh "mongodb://localhost:27017/${1:-test}"
}

# =============================================================================
# Redis
# =============================================================================
if command -v iredis >/dev/null 2>&1; then
    alias rd='iredis'
elif command -v redis-cli >/dev/null 2>&1; then
    alias rd='redis-cli'
fi

if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias rdstart='brew services start redis'
    alias rdstop='brew services stop redis'
    alias rdrestart='brew services restart redis'
    alias rdstatus='brew services info redis'
fi

# Quick connect helpers
rdlocal() {
    if command -v iredis >/dev/null 2>&1; then
        iredis -h localhost -p "${1:-6379}"
    else
        redis-cli -h localhost -p "${1:-6379}"
    fi
}

# =============================================================================
# SQLite
# =============================================================================
alias sq='sqlite3'
alias sqr='sqlite3 -readonly'
alias sqh='sqlite3 -header -column'  # Pretty output

# =============================================================================
# Neo4j
# =============================================================================
if command -v cypher-shell >/dev/null 2>&1; then
    alias neo='cypher-shell'
fi

if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias neostart='brew services start neo4j'
    alias neostop='brew services stop neo4j'
    alias neorestart='brew services restart neo4j'
    alias neostatus='brew services info neo4j'
fi

# Quick connect helpers
neolocal() {
    cypher-shell -u "${1:-neo4j}" -p "${2:-neo4j}" -a "bolt://localhost:7687"
}

# =============================================================================
# Kafka
# =============================================================================
if command -v kafka-console-producer >/dev/null 2>&1; then
    # Producer/Consumer shortcuts
    alias kprod='kafka-console-producer'
    alias kcons='kafka-console-consumer'
    alias ktop='kafka-topics'
    alias kgroups='kafka-consumer-groups'
fi

if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    alias kafkastart='brew services start kafka'
    alias kafkastop='brew services stop kafka'
    alias kafkarestart='brew services restart kafka'
    alias kafkastatus='brew services info kafka'
    alias zkstart='brew services start zookeeper'
    alias zkstop='brew services stop zookeeper'
fi

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
