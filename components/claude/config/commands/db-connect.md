---
description: Help connect to various database systems
argument-hint: [database: postgres|mysql|mongodb|redis|sqlite|neo4j]
allowed-tools: Read, Bash
---

## Task

Help the user connect to a database. Determine which database they need and provide the appropriate connection method.

## Quick Connect Commands

### PostgreSQL
```bash
# Local connection (dotfiles function)
pglocal                    # Connect as postgres to postgres db
pglocal myuser mydb        # Connect as myuser to mydb

# Using pgcli (alias: pg)
pg -h localhost -U postgres -d mydb
pg postgres://user:pass@host:5432/dbname

# Using psql (alias: psqlc)
psqlc -h localhost -U postgres mydb
```

### MySQL
```bash
# Local connection (dotfiles function)
mylocal                    # Connect as root
mylocal myuser mydb        # Connect as myuser to mydb

# Using mycli (alias: my)
my -h localhost -u root
my mysql://user:pass@host:3306/dbname
```

### MongoDB
```bash
# Local connection (dotfiles function)
mongolocal                 # Connect to test database
mongolocal mydb            # Connect to specific database

# Using mongosh (alias: mongo or msh)
mongo mongodb://localhost:27017/mydb
mongo mongodb://user:pass@host:27017/mydb?authSource=admin
```

### Redis
```bash
# Local connection (dotfiles function)
rdlocal                    # Connect to default port 6379
rdlocal 6380               # Connect to custom port

# Using redis-cli (alias: rd)
rd -h localhost -p 6379
rd -h localhost -p 6379 -a password
```

### SQLite
```bash
# Using sqlite3 (alias: sq)
sq mydb.sqlite             # Open database file
sqr mydb.sqlite            # Open read-only
sqh mydb.sqlite            # Open with headers and columns

# Create new database
sq newdb.sqlite "CREATE TABLE test (id INTEGER PRIMARY KEY);"
```

### Neo4j
```bash
# Local connection (dotfiles function)
neolocal                   # Connect as neo4j/neo4j
neolocal myuser mypass     # Custom credentials

# Using cypher-shell (alias: neo)
neo -u neo4j -p password -a bolt://localhost:7687
```

## Connection String Formats

| Database | Format |
|----------|--------|
| PostgreSQL | `postgres://user:pass@host:5432/dbname` |
| MySQL | `mysql://user:pass@host:3306/dbname` |
| MongoDB | `mongodb://user:pass@host:27017/dbname` |
| Redis | `redis://user:pass@host:6379/0` |

## Troubleshooting

### Database Not Running
```bash
# Check status first
db_status

# Start services (macOS)
pgstart    # PostgreSQL
mystart    # MySQL
mongostart # MongoDB
rdstart    # Redis
neostart   # Neo4j
```

### Connection Refused
- Verify database is running: `db_status`
- Check port is correct (default ports above)
- Verify firewall allows connection
- Check authentication credentials

### Authentication Failed
- Reset password through database admin tools
- Check user has appropriate permissions
- Verify authentication method (password, md5, scram, etc.)

## Service Management (macOS)

### Start/Stop Individual Databases
```bash
# PostgreSQL
pgstart / pgstop / pgrestart / pgstatus

# MySQL
mystart / mystop / myrestart / mystatus

# MongoDB
mongostart / mongostop / mongorestart / mongostatus

# Redis
rdstart / rdstop / rdrestart / rdstatus

# Neo4j
neostart / neostop / neorestart / neostatus
```

### Start/Stop All
```bash
db_start_all   # Start PostgreSQL, Redis, MongoDB
db_stop_all    # Stop all
```

## Context

@components/database/functions.sh
@components/database/aliases.sh
