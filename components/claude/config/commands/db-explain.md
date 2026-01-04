---
description: Explain database concepts, queries, and configurations
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about databases. If no specific topic provided, give an overview of the database component.

## Topics

### Database Systems
- **postgresql** - PostgreSQL features, data types, extensions
- **mysql** - MySQL features, storage engines, syntax
- **mongodb** - MongoDB document model, aggregations
- **redis** - Redis data structures, caching patterns
- **sqlite** - SQLite use cases, limitations
- **neo4j** - Graph database concepts, Cypher queries
- **kafka** - Message streaming, topics, partitions

### Concepts
- **indexes** - How indexes work, when to use them
- **transactions** - ACID properties, isolation levels
- **joins** - SQL join types and when to use each
- **normalization** - Database normalization forms
- **replication** - Master-slave, multi-master patterns
- **sharding** - Horizontal scaling strategies

### Operations
- **backup** - Backup strategies for each database
- **restore** - Recovery procedures
- **migration** - Schema migrations, data migrations
- **monitoring** - Performance monitoring tools

## Context

@components/database/component.yaml
@components/database/functions.sh
@components/database/aliases.sh

## Response Format

When explaining a topic:

1. **Definition** - What it is in simple terms
2. **How it works** - Technical details
3. **Examples** - Practical code/query examples
4. **Dotfiles integration** - Available functions/aliases
5. **Best practices** - Common patterns and pitfalls

## Quick Reference

### Dotfiles Functions
- `pglocal [user] [db]` - PostgreSQL connection
- `mylocal [user] [db]` - MySQL connection
- `mongolocal [db]` - MongoDB connection
- `rdlocal [port]` - Redis connection
- `neolocal [user] [pass]` - Neo4j connection
- `db_status` - Check all database services
- `db_start_all` / `db_stop_all` - Service management

### Key Aliases
| Alias | Command | Database |
|-------|---------|----------|
| `pg` | pgcli | PostgreSQL |
| `my` | mycli | MySQL |
| `mongo` | mongosh | MongoDB |
| `rd` | redis-cli | Redis |
| `sq` | sqlite3 | SQLite |
| `neo` | cypher-shell | Neo4j |
