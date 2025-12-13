# Database Connection Examples

## Overview

Database apps in Mist are deployed on the internal Docker network (`traefik-net`). They are NOT exposed to the host machine or internet, ensuring security. Other containers on the same network can connect using the container name as the hostname.

---

## Network Architecture

```
┌─────────────────────────────────────────────────┐
│               traefik-net (Internal)            │
│                                                 │
│  ┌──────────────┐    ┌──────────────┐          │
│  │  Web App     │───▶│  PostgreSQL  │          │
│  │ app-12345    │    │ app-54321    │          │
│  └──────────────┘    └──────────────┘          │
│         │                                       │
│         │            ┌──────────────┐          │
│         └───────────▶│    Redis     │          │
│                      │ app-67890    │          │
│                      └──────────────┘          │
│                                                 │
└─────────────────────────────────────────────────┘
         │
         │ (Only web apps get external access)
         ▼
  ┌─────────────┐
  │  Internet   │
  │   Users     │
  └─────────────┘
```

**Key Points:**
- ✅ Internal communication works via container names
- ❌ External access blocked for database/service apps
- ✅ Secure by default

---

## Connection String Format

### General Format
```
protocol://username:password@container_name:port/database
```

**Important:** Use the **container name** (e.g., `app-522001`) as the hostname, NOT `localhost` or an IP address.

---

## PostgreSQL Connections

### Connection String Examples

**Basic Connection:**
```bash
postgresql://postgres:iUXK91^&RVvmWXal@app-522001:5432/postgres
```

**With Database Name:**
```bash
postgresql://postgres:YOUR_PASSWORD@app-522001:5432/myapp
```

**With Additional Parameters:**
```bash
postgresql://postgres:YOUR_PASSWORD@app-522001:5432/myapp?sslmode=disable
```

### Environment Variables for Web Apps

Add these to your web app's environment variables:

```env
DATABASE_URL=postgresql://postgres:iUXK91^&RVvmWXal@app-522001:5432/postgres
DB_HOST=app-522001
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=iUXK91^&RVvmWXal
DB_NAME=postgres
```

### Code Examples

**Node.js (pg library):**
```javascript
const { Pool } = require('pg');

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  // OR
  host: 'app-522001',
  port: 5432,
  user: 'postgres',
  password: process.env.DB_PASSWORD,
  database: 'postgres'
});

// Test connection
pool.query('SELECT NOW()', (err, res) => {
  console.log(err, res);
  pool.end();
});
```

**Python (psycopg2):**
```python
import psycopg2
import os

conn = psycopg2.connect(os.environ['DATABASE_URL'])
# OR
conn = psycopg2.connect(
    host='app-522001',
    port=5432,
    user='postgres',
    password=os.environ['DB_PASSWORD'],
    dbname='postgres'
)

cur = conn.cursor()
cur.execute('SELECT version()')
print(cur.fetchone())
```

**Go:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "os"
)

connStr := os.Getenv("DATABASE_URL")
db, err := sql.Open("postgres", connStr)
// OR
connStr := fmt.Sprintf("host=app-522001 port=5432 user=postgres "+
    "password=%s dbname=postgres sslmode=disable",
    os.Getenv("DB_PASSWORD"))
```

**Java (JDBC):**
```java
String url = System.getenv("DATABASE_URL");
// OR
String url = "jdbc:postgresql://app-522001:5432/postgres";
Properties props = new Properties();
props.setProperty("user", "postgres");
props.setProperty("password", System.getenv("DB_PASSWORD"));
Connection conn = DriverManager.getConnection(url, props);
```

---

## Redis Connections

### Connection Examples

**Basic Connection (no password):**
```bash
redis://app-243241:6379
```

**With Password:**
```bash
redis://:jgU1VI5T5rTdXNFe@app-243241:6379
```

**Redis CLI:**
```bash
redis-cli -h app-243241 -p 6379
# With password
redis-cli -h app-243241 -p 6379 -a jgU1VI5T5rTdXNFe
```

### Environment Variables

```env
REDIS_URL=redis://:jgU1VI5T5rTdXNFe@app-243241:6379
REDIS_HOST=app-243241
REDIS_PORT=6379
REDIS_PASSWORD=jgU1VI5T5rTdXNFe
```

### Code Examples

**Node.js (ioredis):**
```javascript
const Redis = require('ioredis');

const redis = new Redis(process.env.REDIS_URL);
// OR
const redis = new Redis({
  host: 'app-243241',
  port: 6379,
  password: process.env.REDIS_PASSWORD
});

redis.set('key', 'value');
redis.get('key', (err, result) => {
  console.log(result); // 'value'
});
```

**Python (redis-py):**
```python
import redis
import os

r = redis.from_url(os.environ['REDIS_URL'])
# OR
r = redis.Redis(
    host='app-243241',
    port=6379,
    password=os.environ['REDIS_PASSWORD'],
    decode_responses=True
)

r.set('key', 'value')
print(r.get('key'))  # 'value'
```

**Go (go-redis):**
```go
import (
    "github.com/go-redis/redis/v8"
    "context"
)

rdb := redis.NewClient(&redis.Options{
    Addr:     "app-243241:6379",
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       0,
})

ctx := context.Background()
rdb.Set(ctx, "key", "value", 0)
```

---

## MySQL/MariaDB Connections

### Connection String Examples

**MySQL:**
```bash
mysql://root:YOUR_PASSWORD@app-xxxxx:3306/mydb
```

**MariaDB:**
```bash
mariadb://root:YOUR_PASSWORD@app-xxxxx:3306/mydb
```

### Environment Variables

```env
DATABASE_URL=mysql://root:Kp9@xL2#mN8@app-xxxxx:3306/mydb
MYSQL_HOST=app-xxxxx
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=Kp9@xL2#mN8
MYSQL_DATABASE=mydb
```

### Code Examples

**Node.js (mysql2):**
```javascript
const mysql = require('mysql2/promise');

const connection = await mysql.createConnection({
  host: 'app-xxxxx',
  port: 3306,
  user: 'root',
  password: process.env.MYSQL_PASSWORD,
  database: 'mydb'
});

const [rows] = await connection.execute('SELECT * FROM users');
```

**Python (PyMySQL):**
```python
import pymysql

connection = pymysql.connect(
    host='app-xxxxx',
    port=3306,
    user='root',
    password=os.environ['MYSQL_PASSWORD'],
    database='mydb'
)

with connection.cursor() as cursor:
    cursor.execute('SELECT VERSION()')
    print(cursor.fetchone())
```

---

## MongoDB Connections

### Connection String Examples

**Basic:**
```bash
mongodb://admin:YOUR_PASSWORD@app-xxxxx:27017
```

**With Database:**
```bash
mongodb://admin:YOUR_PASSWORD@app-xxxxx:27017/mydb
```

**With Auth Database:**
```bash
mongodb://admin:YOUR_PASSWORD@app-xxxxx:27017/mydb?authSource=admin
```

### Environment Variables

```env
MONGODB_URI=mongodb://admin:Xy7!mK9@pL3#@app-xxxxx:27017/mydb?authSource=admin
MONGO_HOST=app-xxxxx
MONGO_PORT=27017
MONGO_USER=admin
MONGO_PASSWORD=Xy7!mK9@pL3#
MONGO_DATABASE=mydb
```

### Code Examples

**Node.js (mongodb driver):**
```javascript
const { MongoClient } = require('mongodb');

const client = new MongoClient(process.env.MONGODB_URI);
await client.connect();

const db = client.db('mydb');
const collection = db.collection('users');
```

**Python (pymongo):**
```python
from pymongo import MongoClient
import os

client = MongoClient(os.environ['MONGODB_URI'])
db = client['mydb']
collection = db['users']
```

---

## Testing Connections

### From Another Container

You can test database connections by running temporary containers:

**Test PostgreSQL:**
```bash
docker run --rm --network traefik-net postgres:16-alpine \
  psql "postgresql://postgres:PASSWORD@app-522001:5432/postgres" \
  -c "SELECT version();"
```

**Test Redis:**
```bash
docker run --rm --network traefik-net redis:7-alpine \
  redis-cli -h app-243241 -p 6379 PING
```

**Test MySQL:**
```bash
docker run --rm --network traefik-net mysql:8 \
  mysql -h app-xxxxx -u root -pPASSWORD -e "SELECT VERSION();"
```

**Test MongoDB:**
```bash
docker run --rm --network traefik-net mongo:7 \
  mongosh "mongodb://admin:PASSWORD@app-xxxxx:27017/mydb?authSource=admin" \
  --eval "db.version()"
```

### Verify External Access is Blocked

These commands should FAIL (expected behavior):

```bash
# Try to connect from host machine
psql -h localhost -p 5432 -U postgres
# Result: Connection refused ✅

nc -zv localhost 5432
# Result: Connection refused ✅

redis-cli -h localhost -p 6379
# Result: Connection refused ✅
```

---

## Common Issues & Solutions

### Issue: "Could not translate host name to address"

**Error:**
```
could not translate host name "app-522001" to address: Name or service not known
```

**Cause:** Your web app container is not on the `traefik-net` network.

**Solution:**
```bash
# Check container network
docker inspect your-web-app | grep -A 5 Networks

# If not on traefik-net, stop and redeploy via Mist UI
```

### Issue: "Connection refused"

**Error:**
```
could not connect to server: Connection refused
    Is the server running on host "app-522001" and accepting TCP/IP connections on port 5432?
```

**Cause:** Database container is not running or wrong container name.

**Solution:**
```bash
# Check if database is running
docker ps | grep app-522001

# Check if it's on the correct network
docker network inspect traefik-net | grep app-522001
```

### Issue: "Authentication failed"

**Error:**
```
password authentication failed for user "postgres"
```

**Cause:** Wrong password in connection string.

**Solution:**
1. Check environment variables in Mist UI (App → Environment tab)
2. Find the correct `POSTGRES_PASSWORD` value
3. Update your connection string

### Issue: Using localhost instead of container name

**Wrong:**
```javascript
host: 'localhost',  // ❌ Won't work
host: '127.0.0.1',  // ❌ Won't work
```

**Correct:**
```javascript
host: 'app-522001',  // ✅ Use container name
```

---

## Best Practices

### 1. Use Environment Variables

**Don't hardcode:**
```javascript
const db = new Pool({
  host: 'app-522001',           // ❌ Hardcoded
  password: 'iUXK91^&RVvmWXal' // ❌ Hardcoded in code
});
```

**Use env vars:**
```javascript
const db = new Pool({
  connectionString: process.env.DATABASE_URL  // ✅ From environment
});
```

### 2. Store Connection Strings in Mist

Add these environment variables to your web app in Mist:

```env
DATABASE_URL=postgresql://postgres:PASSWORD@app-522001:5432/postgres
REDIS_URL=redis://:PASSWORD@app-243241:6379
```

### 3. Handle Connection Errors

```javascript
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,
  connectionTimeoutMillis: 2000,
  idleTimeoutMillis: 30000
});

pool.on('error', (err, client) => {
  console.error('Unexpected database error:', err);
});
```

### 4. Use Connection Pooling

**Node.js:**
```javascript
const { Pool } = require('pg');
const pool = new Pool({ connectionString: process.env.DATABASE_URL });
```

**Python:**
```python
from sqlalchemy import create_engine

engine = create_engine(
    os.environ['DATABASE_URL'],
    pool_size=10,
    max_overflow=20
)
```

### 5. Test Connections on Startup

```javascript
async function testDatabaseConnection() {
  try {
    const client = await pool.connect();
    await client.query('SELECT NOW()');
    client.release();
    console.log('✅ Database connected successfully');
  } catch (err) {
    console.error('❌ Database connection failed:', err);
    process.exit(1);
  }
}

testDatabaseConnection();
```

---

## Example: Full Stack Application

### Scenario: Node.js API + PostgreSQL + Redis

**1. Create PostgreSQL database in Mist:**
- Name: `my-postgres`
- Container: `app-100001`
- Password: `dbpass123`

**2. Create Redis in Mist:**
- Name: `my-redis`
- Container: `app-100002`
- Password: `redispass456`

**3. Create Web App in Mist:**
- Name: `my-api`
- Port: 3000
- Add environment variables:

```env
DATABASE_URL=postgresql://postgres:dbpass123@app-100001:5432/postgres
REDIS_URL=redis://:redispass456@app-100002:6379
PORT=3000
NODE_ENV=production
```

**4. Application Code (`server.js`):**

```javascript
const express = require('express');
const { Pool } = require('pg');
const Redis = require('ioredis');

const app = express();
const port = process.env.PORT || 3000;

// Connect to PostgreSQL
const db = new Pool({
  connectionString: process.env.DATABASE_URL
});

// Connect to Redis
const redis = new Redis(process.env.REDIS_URL);

// Health check endpoint
app.get('/health', async (req, res) => {
  try {
    // Test database
    await db.query('SELECT NOW()');
    
    // Test Redis
    await redis.ping();
    
    res.json({ 
      status: 'healthy',
      database: 'connected',
      cache: 'connected'
    });
  } catch (err) {
    res.status(500).json({ 
      status: 'unhealthy',
      error: err.message 
    });
  }
});

// Example endpoint with caching
app.get('/users/:id', async (req, res) => {
  const { id } = req.params;
  
  // Check cache
  const cached = await redis.get(`user:${id}`);
  if (cached) {
    return res.json(JSON.parse(cached));
  }
  
  // Query database
  const result = await db.query(
    'SELECT * FROM users WHERE id = $1',
    [id]
  );
  
  if (result.rows.length === 0) {
    return res.status(404).json({ error: 'User not found' });
  }
  
  const user = result.rows[0];
  
  // Cache for 5 minutes
  await redis.setex(`user:${id}`, 300, JSON.stringify(user));
  
  res.json(user);
});

app.listen(port, () => {
  console.log(`API running on port ${port}`);
});
```

**5. Deploy:**
- Push code to Git
- Deploy via Mist UI
- Application will connect to both databases automatically

---

## Security Notes

### ✅ Secure by Default

- Database containers are NOT exposed to the internet
- Only accessible from containers on `traefik-net`
- Passwords stored as environment variables
- No direct host port mapping

### 🔒 Additional Security Recommendations

1. **Rotate passwords regularly**
2. **Use strong passwords** (16+ characters, mixed case, symbols)
3. **Create separate database users** for each application
4. **Use read-only users** where possible
5. **Enable SSL/TLS** for production databases
6. **Regular backups** (implement volume persistence)

---

## Conclusion

Database connections in Mist work seamlessly using container names as hostnames. The internal network (`traefik-net`) provides secure, isolated communication between your applications and databases without exposing sensitive services to the internet.

**Key Takeaway:** Use the container name (e.g., `app-522001`) in your connection strings, and store credentials as environment variables in Mist.
