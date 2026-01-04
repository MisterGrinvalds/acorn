---
description: Help with database backup and restore operations
argument-hint: [database: postgres|mysql|mongodb|redis|sqlite]
allowed-tools: Read, Bash
---

## Task

Help the user backup or restore databases. Provide appropriate commands based on the database system.

## PostgreSQL Backup

### Full Database Backup
```bash
# Single database
pg_dump mydb > mydb_backup.sql
pg_dump -Fc mydb > mydb_backup.dump  # Custom format (compressed)

# With connection details
pg_dump -h localhost -U postgres mydb > backup.sql

# All databases
pg_dumpall > all_databases.sql
pg_dumpall -h localhost -U postgres > all_databases.sql
```

### Specific Tables
```bash
# Single table
pg_dump -t users mydb > users_backup.sql

# Multiple tables
pg_dump -t users -t orders mydb > tables_backup.sql

# Schema only (no data)
pg_dump --schema-only mydb > schema.sql

# Data only
pg_dump --data-only mydb > data.sql
```

### Restore PostgreSQL
```bash
# SQL format
psql mydb < backup.sql

# Custom format
pg_restore -d mydb mydb_backup.dump

# Create database during restore
pg_restore -C -d postgres mydb_backup.dump
```

## MySQL Backup

### Full Database Backup
```bash
# Single database
mysqldump mydb > mydb_backup.sql
mysqldump -u root -p mydb > mydb_backup.sql

# All databases
mysqldump --all-databases > all_databases.sql

# With routines and events
mysqldump --routines --events mydb > backup.sql
```

### Specific Tables
```bash
# Single table
mysqldump mydb users > users_backup.sql

# Multiple tables
mysqldump mydb users orders > tables_backup.sql
```

### Restore MySQL
```bash
mysql mydb < backup.sql
mysql -u root -p mydb < backup.sql

# Create database if not exists
mysql -e "CREATE DATABASE IF NOT EXISTS mydb"
mysql mydb < backup.sql
```

## MongoDB Backup

### Full Backup
```bash
# Entire server
mongodump --out ./backup

# Single database
mongodump --db mydb --out ./backup

# With authentication
mongodump --uri "mongodb://user:pass@localhost:27017/mydb" --out ./backup
```

### Single Collection
```bash
mongodump --db mydb --collection users --out ./backup
```

### Export to JSON
```bash
# Single collection
mongoexport --db mydb --collection users --out users.json

# As JSON array
mongoexport --db mydb --collection users --jsonArray --out users.json
```

### Restore MongoDB
```bash
# Full restore
mongorestore ./backup

# Single database
mongorestore --db mydb ./backup/mydb

# Drop existing before restore
mongorestore --drop ./backup
```

## Redis Backup

### RDB Snapshot
```bash
# Trigger manual save
redis-cli SAVE        # Blocking
redis-cli BGSAVE      # Background

# Check last save time
redis-cli LASTSAVE

# Find RDB file location
redis-cli CONFIG GET dir
redis-cli CONFIG GET dbfilename
# Usually: /var/lib/redis/dump.rdb or /opt/homebrew/var/db/redis/dump.rdb
```

### Copy Backup
```bash
# Copy RDB file while Redis is running
redis-cli BGSAVE
cp /opt/homebrew/var/db/redis/dump.rdb ./redis_backup.rdb
```

### AOF Backup
```bash
# Enable AOF persistence
redis-cli CONFIG SET appendonly yes

# Rewrite AOF file
redis-cli BGREWRITEAOF
```

### Restore Redis
```bash
# Stop Redis
rdstop

# Replace dump.rdb
cp redis_backup.rdb /opt/homebrew/var/db/redis/dump.rdb

# Start Redis
rdstart
```

## SQLite Backup

### Simple Copy
```bash
# SQLite is a single file - just copy it
cp mydb.sqlite mydb_backup.sqlite

# With date stamp
cp mydb.sqlite "mydb_$(date +%Y%m%d_%H%M%S).sqlite"
```

### Using .backup Command
```bash
sqlite3 mydb.sqlite ".backup backup.sqlite"
```

### Export to SQL
```bash
sqlite3 mydb.sqlite .dump > backup.sql
```

### Restore SQLite
```bash
# From SQL dump
sqlite3 newdb.sqlite < backup.sql

# From backup file
cp backup.sqlite mydb.sqlite
```

## Automated Backup Script

```bash
#!/bin/bash
# backup_databases.sh

BACKUP_DIR="$HOME/backups/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

# PostgreSQL
pg_dump mydb > "$BACKUP_DIR/postgres_mydb.sql"

# MySQL
mysqldump mydb > "$BACKUP_DIR/mysql_mydb.sql"

# MongoDB
mongodump --db mydb --out "$BACKUP_DIR/mongodb"

# Redis
redis-cli BGSAVE
sleep 2
cp /opt/homebrew/var/db/redis/dump.rdb "$BACKUP_DIR/redis.rdb"

echo "Backups saved to $BACKUP_DIR"
```

## Best Practices

1. **Regular backups** - Daily for active databases
2. **Test restores** - Periodically verify backups work
3. **Offsite storage** - Copy to cloud storage
4. **Retention policy** - Keep backups for appropriate time
5. **Document process** - Record backup and restore steps
