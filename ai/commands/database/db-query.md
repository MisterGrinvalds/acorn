---
description: Help write and optimize database queries
argument-hint: [database: postgres|mysql|mongodb|redis|neo4j]
allowed-tools: Read, Bash
---

## Task

Help the user write, understand, or optimize database queries. Provide syntax-appropriate examples for the target database.

## PostgreSQL Queries

### Basic CRUD
```sql
-- Select
SELECT * FROM users WHERE active = true;
SELECT id, name, email FROM users ORDER BY created_at DESC LIMIT 10;

-- Insert
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com');
INSERT INTO users (name, email) VALUES
  ('Alice', 'alice@example.com'),
  ('Bob', 'bob@example.com');

-- Update
UPDATE users SET active = false WHERE last_login < NOW() - INTERVAL '1 year';

-- Delete
DELETE FROM users WHERE id = 123;
```

### Joins
```sql
-- Inner join
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- Left join (all users, even without orders)
SELECT u.name, COALESCE(COUNT(o.id), 0) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id;

-- Multiple joins
SELECT u.name, p.name as product, o.quantity
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN products p ON o.product_id = p.id;
```

### Aggregations
```sql
-- Group by with having
SELECT category, COUNT(*) as count, AVG(price) as avg_price
FROM products
GROUP BY category
HAVING COUNT(*) > 5;

-- Window functions
SELECT name, department, salary,
  RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank,
  AVG(salary) OVER (PARTITION BY department) as dept_avg
FROM employees;
```

### Query Optimization
```sql
-- Analyze query plan
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';

-- Create index
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at);

-- Partial index
CREATE INDEX idx_active_users ON users(email) WHERE active = true;
```

## MySQL Queries

### Basic CRUD
```sql
-- Similar to PostgreSQL with some differences

-- Insert with ON DUPLICATE KEY
INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- REPLACE (delete and insert)
REPLACE INTO users (id, name) VALUES (1, 'Alice');
```

### MySQL-Specific
```sql
-- LIMIT with OFFSET
SELECT * FROM users LIMIT 10 OFFSET 20;
-- or
SELECT * FROM users LIMIT 20, 10;

-- Full-text search
SELECT * FROM articles
WHERE MATCH(title, body) AGAINST('database' IN NATURAL LANGUAGE MODE);

-- JSON queries
SELECT * FROM users WHERE JSON_EXTRACT(settings, '$.theme') = 'dark';
```

## MongoDB Queries

### Basic CRUD
```javascript
// Find
db.users.find({ active: true })
db.users.find({ age: { $gte: 18 } }).sort({ name: 1 }).limit(10)

// Insert
db.users.insertOne({ name: "Alice", email: "alice@example.com" })
db.users.insertMany([
  { name: "Alice", email: "alice@example.com" },
  { name: "Bob", email: "bob@example.com" }
])

// Update
db.users.updateOne({ _id: ObjectId("...") }, { $set: { active: false } })
db.users.updateMany({ lastLogin: { $lt: new Date("2023-01-01") } }, { $set: { active: false } })

// Delete
db.users.deleteOne({ _id: ObjectId("...") })
db.users.deleteMany({ active: false })
```

### Aggregation Pipeline
```javascript
// Count by status
db.orders.aggregate([
  { $group: { _id: "$status", count: { $sum: 1 } } }
])

// Complex aggregation
db.orders.aggregate([
  { $match: { status: "completed" } },
  { $group: {
      _id: "$customerId",
      totalSpent: { $sum: "$total" },
      orderCount: { $sum: 1 }
  }},
  { $sort: { totalSpent: -1 } },
  { $limit: 10 }
])

// Lookup (join)
db.orders.aggregate([
  { $lookup: {
      from: "users",
      localField: "userId",
      foreignField: "_id",
      as: "user"
  }},
  { $unwind: "$user" }
])
```

### Indexes
```javascript
// Create index
db.users.createIndex({ email: 1 }, { unique: true })
db.orders.createIndex({ userId: 1, createdAt: -1 })

// Text index
db.articles.createIndex({ title: "text", body: "text" })
db.articles.find({ $text: { $search: "database" } })
```

## Redis Commands

### String Operations
```bash
SET key "value"
GET key
MSET key1 "v1" key2 "v2"
MGET key1 key2
INCR counter
INCRBY counter 5
EXPIRE key 3600  # TTL in seconds
TTL key
```

### List Operations
```bash
LPUSH queue "item1" "item2"
RPUSH queue "item3"
LPOP queue
RPOP queue
LRANGE queue 0 -1  # All items
LLEN queue
```

### Set Operations
```bash
SADD tags "redis" "database"
SMEMBERS tags
SISMEMBER tags "redis"
SINTER tags1 tags2
SUNION tags1 tags2
```

### Hash Operations
```bash
HSET user:1 name "Alice" email "alice@example.com"
HGET user:1 name
HGETALL user:1
HINCRBY user:1 visits 1
```

### Sorted Set Operations
```bash
ZADD leaderboard 100 "player1" 200 "player2"
ZRANGE leaderboard 0 -1 WITHSCORES
ZREVRANGE leaderboard 0 9 WITHSCORES  # Top 10
ZRANK leaderboard "player1"
```

## Neo4j Cypher Queries

### Basic CRUD
```cypher
// Create
CREATE (p:Person {name: 'Alice', age: 30})
CREATE (a:Person {name: 'Alice'})-[:KNOWS]->(b:Person {name: 'Bob'})

// Read
MATCH (p:Person) RETURN p
MATCH (p:Person {name: 'Alice'}) RETURN p
MATCH (a:Person)-[:KNOWS]->(b:Person) RETURN a.name, b.name

// Update
MATCH (p:Person {name: 'Alice'}) SET p.age = 31

// Delete
MATCH (p:Person {name: 'Alice'}) DELETE p
MATCH (p:Person {name: 'Alice'}) DETACH DELETE p  # With relationships
```

### Graph Traversal
```cypher
// Friends of friends
MATCH (a:Person {name: 'Alice'})-[:KNOWS*2]->(fof)
RETURN DISTINCT fof.name

// Shortest path
MATCH path = shortestPath((a:Person {name: 'Alice'})-[*]-(b:Person {name: 'Bob'}))
RETURN path

// All paths up to 4 hops
MATCH path = (a:Person {name: 'Alice'})-[*1..4]-(b:Person {name: 'Bob'})
RETURN path
```

## Connection Commands

```bash
# PostgreSQL
pglocal
pg -h localhost -U postgres -d mydb -c "SELECT * FROM users;"

# MySQL
mylocal
my -h localhost -u root -e "SELECT * FROM users;"

# MongoDB
mongolocal
mongo --eval "db.users.find()"

# Redis
rdlocal
rd GET key

# Neo4j
neolocal
neo -u neo4j -p password "MATCH (n) RETURN n LIMIT 10"
```
