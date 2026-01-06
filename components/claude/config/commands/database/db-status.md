---
description: Check status of database services
argument-hint: [database: all|postgres|mysql|mongodb|redis|neo4j]
allowed-tools: Bash
---

## Task

Check the status of database services and report which are running, stopped, or not installed.

## Quick Status Check

```bash
# Check all databases at once (dotfiles function)
db_status
```

This shows status for:
- PostgreSQL
- MySQL
- MongoDB
- Redis
- Neo4j

## Individual Status Checks

### PostgreSQL
```bash
# Dotfiles alias
pgstatus

# Manual check
pg_isready
# Returns 0 if accepting connections

# Detailed status (macOS)
brew services info postgresql@14
```

### MySQL
```bash
# Dotfiles alias
mystatus

# Manual check
mysqladmin ping -u root

# Detailed status (macOS)
brew services info mysql
```

### MongoDB
```bash
# Dotfiles alias
mongostatus

# Manual check
mongosh --eval "db.runCommand({ping:1})" --quiet

# Detailed status (macOS)
brew services info mongodb-community
```

### Redis
```bash
# Dotfiles alias
rdstatus

# Manual check
redis-cli ping
# Returns PONG if running

# Detailed status (macOS)
brew services info redis
```

### Neo4j
```bash
# Dotfiles alias
neostatus

# Manual check
neo4j status

# Detailed status (macOS)
brew services info neo4j
```

### Kafka
```bash
# Dotfiles alias
kafkastatus

# Detailed status (macOS)
brew services info kafka
brew services info zookeeper
```

## Service Management

### Start Services
```bash
# Individual
pgstart      # PostgreSQL
mystart      # MySQL
mongostart   # MongoDB
rdstart      # Redis
neostart     # Neo4j
kafkastart   # Kafka

# All common databases
db_start_all
```

### Stop Services
```bash
# Individual
pgstop       # PostgreSQL
mystop       # MySQL
mongostop    # MongoDB
rdstop       # Redis
neostop      # Neo4j
kafkastop    # Kafka

# All
db_stop_all
```

### Restart Services
```bash
pgrestart    # PostgreSQL
myrestart    # MySQL
mongorestart # MongoDB
rdrestart    # Redis
neorestart   # Neo4j
kafkarestart # Kafka
```

## Health Checks

### PostgreSQL Health
```bash
# Connection count
psql -c "SELECT count(*) FROM pg_stat_activity;"

# Database sizes
psql -c "SELECT pg_database.datname, pg_size_pretty(pg_database_size(pg_database.datname)) FROM pg_database;"
```

### Redis Health
```bash
redis-cli info stats
redis-cli info memory
redis-cli dbsize
```

### MongoDB Health
```bash
mongosh --eval "db.serverStatus().connections"
mongosh --eval "db.stats()"
```

## Common Issues

### Service Won't Start
```bash
# Check logs (macOS)
cat /opt/homebrew/var/log/postgresql@14.log
cat /opt/homebrew/var/log/redis.log
cat /opt/homebrew/var/log/mongodb/mongo.log

# Check port conflicts
lsof -i :5432  # PostgreSQL
lsof -i :3306  # MySQL
lsof -i :27017 # MongoDB
lsof -i :6379  # Redis
lsof -i :7687  # Neo4j
```

### Permission Issues
```bash
# Fix PostgreSQL data directory
chmod 700 /opt/homebrew/var/postgresql@14

# Fix Redis
chmod 755 /opt/homebrew/var/db/redis
```
