---
description: Interactive coaching session to learn database management
argument-hint: [database: postgres|mysql|mongodb|redis|sqlite|neo4j]
allowed-tools: Read, Glob, Grep, Bash
---

## Task

Guide the user through learning database management interactively.

## Approach

1. **Assess level** - Ask about database experience
2. **Choose database** - Focus on one or provide overview
3. **Progressive exercises** - Start with connections, build to queries
4. **Real-time practice** - Have them run actual commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- What is a database?
- SQL vs NoSQL differences
- Installing database tools
- Connecting to databases
- Basic CRUD operations

### Intermediate
- Complex queries (joins, subqueries)
- Indexing strategies
- Transaction management
- Backup and restore
- Using dotfiles functions

### Advanced
- Query optimization
- Replication setup
- Sharding strategies
- Performance tuning
- Cross-database operations

## Interactive Exercises

### Beginner Exercises
```bash
# Exercise 1: Check installed databases
db_status

# Exercise 2: Connect to PostgreSQL
pglocal  # or: pg -h localhost -U postgres

# Exercise 3: Basic SQL
CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100));
INSERT INTO users (name) VALUES ('Alice'), ('Bob');
SELECT * FROM users;

# Exercise 4: Connect to Redis
rdlocal
SET greeting "Hello"
GET greeting
```

### Intermediate Exercises
```bash
# Exercise 5: PostgreSQL with joins
SELECT u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.name;

# Exercise 6: MongoDB aggregation
mongolocal
db.orders.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])

# Exercise 7: Redis data structures
LPUSH queue "task1" "task2" "task3"
RPOP queue
LRANGE queue 0 -1
```

### Advanced Exercises
```bash
# Exercise 8: PostgreSQL EXPLAIN
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';

# Exercise 9: Create index
CREATE INDEX idx_users_email ON users(email);

# Exercise 10: Neo4j graph query
neolocal
MATCH (a:Person)-[:KNOWS]->(b:Person)
RETURN a.name, b.name;
```

## Database-Specific Tracks

### PostgreSQL Track
1. Connection with pgcli
2. Basic SQL operations
3. Data types and constraints
4. Joins and subqueries
5. Indexes and EXPLAIN
6. Transactions
7. pg_dump backup

### MongoDB Track
1. Connection with mongosh
2. Document CRUD
3. Query operators
4. Aggregation pipeline
5. Indexes
6. mongodump backup

### Redis Track
1. Connection with redis-cli
2. String operations
3. List, Set, Hash operations
4. Pub/Sub basics
5. Expiration and TTL
6. Persistence options

## Context

@components/database/functions.sh
@components/database/aliases.sh

## Coaching Style

- Start with status check (`db_status`)
- Use local connections first
- Progress from simple to complex queries
- Show dotfiles shortcuts for efficiency
- Emphasize backup importance
- Build toward real application patterns
