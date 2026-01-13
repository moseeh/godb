# godb - A Minimal In-Memory Relational Database

A minimal, in-memory relational database implemented in Go, demonstrating clean architecture, constraint enforcement, indexing, and JOIN operations.

## Features

- **Table Creation** with schema definitions (INT, STRING, BOOL types)
- **Constraints**: Primary keys, unique constraints, and NOT NULL enforcement
- **CRUD Operations**: INSERT, SELECT, UPDATE, DELETE with WHERE clauses
- **Hash-based Indexing** for efficient equality lookups
- **INNER JOIN** support with index optimization
- **Two Interfaces**:
  - Interactive REPL for manual database interaction
  - HTTP REST API for programmatic access

## Architecture

godb follows strict separation of concerns:

```
┌─────────────┐     ┌─────────────┐
│    REPL     │     │  Web Server │
│  (CLI)      │     │  (HTTP API) │
└──────┬──────┘     └──────┬──────┘
       │                   │
       └──────────┬────────┘
                  │
         ┌────────▼────────┐
         │  Database Core  │
         │    (Engine)     │
         └─────────────────┘
```

### Core Components

- **engine/**: Database core - tables, rows, constraints, indexes, CRUD, joins
- **parser/**: SQL-like command parsing (no external dependencies)
- **repl/**: Interactive command-line interface
- **web/**: HTTP REST API server
- **cmd/**: Entry points for REPL and web server

## Installation & Usage

### Prerequisites

- Go 1.16 or higher

### Building

```bash
# Build REPL
go build -o godb-repl cmd/repl/main.go

# Build Web Server
go build -o godb-web cmd/web/main.go
```

### Running the REPL

```bash
./godb-repl
```

#### REPL Examples

```sql
-- Create a table
CREATE TABLE users (id INT PRIMARY KEY, name STRING NOT NULL, email STRING UNIQUE)

-- Insert data
INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com')
INSERT INTO users (id, name, email) VALUES (2, 'Bob', 'bob@example.com')

-- Query data
SELECT * FROM users
SELECT name, email FROM users WHERE id = 1

-- Update data
UPDATE users SET name = 'Alice Smith' WHERE id = 1

-- Delete data
DELETE FROM users WHERE id = 2

-- Create another table
CREATE TABLE posts (id INT PRIMARY KEY, user_id INT NOT NULL, title STRING, body STRING)

-- Insert posts
INSERT INTO posts (id, user_id, title, body) VALUES (1, 1, 'First Post', 'Hello World')
INSERT INTO posts (id, user_id, title, body) VALUES (2, 1, 'Second Post', 'Another post')

-- Perform JOIN
SELECT * FROM posts INNER JOIN users ON posts.user_id = users.id
```

### Running the Web Server

```bash
./godb-web
```

The server starts on `http://localhost:8080` with these endpoints:

#### API Endpoints

**Create User**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "name": "Alice", "email": "alice@example.com"}'
```

**Get All Users**
```bash
curl http://localhost:8080/users
```

**Create Post**
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "user_id": 1, "title": "My Post", "body": "Post content"}'
```

**Get Posts with User Data (JOIN)**
```bash
curl http://localhost:8080/posts
```

## Design Decisions

### 1. In-Memory Storage
- All data stored in memory using Go slices and maps
- Fast access but no persistence
- Suitable for demonstrations, testing, and temporary data

### 2. Hash-Based Indexing
- Indexes use `map[interface{}][]int` structure
- O(1) average lookup time for equality conditions
- Automatically created for PRIMARY KEY and UNIQUE columns
- Can be manually created on any column

### 3. Constraint Enforcement
- **Primary Key**: Uniqueness + NOT NULL automatically enforced
- **Unique**: Prevents duplicate values using indexes
- **NOT NULL**: Validates presence during INSERT/UPDATE
- Validation happens before data modification

### 4. Simple SQL Parser
- No external parser libraries
- Hand-written tokenizer and recursive descent parser
- Supports essential SQL subset only
- Trade-off: Limited SQL features for zero dependencies

### 5. Thread-Safe Database
- RWMutex on database-level operations
- Safe for concurrent reads, exclusive writes
- Table-level locking (not row-level)

### 6. Join Implementation
- Nested loop join algorithm
- Optimizes right table lookup using index if available
- Only INNER JOIN with equality condition supported
- Column names prefixed with table names (e.g., `users.id`)

## Project Structure

```
godb/
├── cmd/
│   ├── repl/main.go          # REPL entry point
│   └── web/main.go           # Web server entry point
├── engine/
│   ├── database.go           # Database and table registry
│   ├── table.go              # Table schema and storage
│   ├── row.go                # Row representation
│   ├── constraints.go        # Constraint validation
│   ├── index.go              # Hash-based indexes
│   ├── crud.go               # CRUD operations
│   ├── join.go               # INNER JOIN logic
│   └── errors.go             # Domain errors
├── parser/
│   ├── ast.go                # Command structures
│   ├── tokenizer.go          # Input tokenization
│   └── parser.go             # Command parsing
├── repl/
│   ├── repl.go               # REPL loop
│   └── printer.go            # Output formatting
├── web/
│   ├── server.go             # HTTP server setup
│   ├── handlers.go           # Route handlers
│   └── dto.go                # Request/response models
├── tests/
│   ├── engine/               # Engine tests
│   └── parser/               # Parser tests
├── go.mod
└── README.md
```

## Testing

Run all tests:

```bash
# Run all tests
go test ./tests/...

# Run with verbose output
go test ./tests/... -v

# Run specific test suite
go test ./tests/engine/...
go test ./tests/parser/...
```

Test coverage includes:
- INSERT with constraint violations
- SELECT with and without indexes
- UPDATE with constraint revalidation
- DELETE operations
- INNER JOIN with and without indexes
- Parser correctness for all command types

## Limitations (By Design)

The following are intentionally **not implemented**:

- **Persistence**: No disk storage or WAL
- **Transactions**: No ACID guarantees, rollback, or commit
- **Advanced SQL**: No GROUP BY, ORDER BY, subqueries, or aggregations
- **Query Optimization**: No query planner or cost-based optimization
- **Authentication**: No user management or access control
- **Network Protocol**: Web server uses HTTP/JSON, not a database protocol
- **Data Types**: Limited to INT, STRING, BOOL

These limitations are deliberate to maintain simplicity and focus on core database concepts.

## Performance Characteristics

- **INSERT**: O(1) with indexing overhead
- **SELECT with indexed equality**: O(1) average
- **SELECT with scan**: O(n)
- **UPDATE/DELETE**: O(n) for condition evaluation
- **JOIN with index**: O(n) for left table, O(1) per right lookup
- **JOIN without index**: O(n * m) nested loop

## Dependencies

- **Go Standard Library Only**
  - `net/http` for web server
  - `encoding/json` for API serialization
  - `bufio` for REPL input
  - `testing` for tests

No external dependencies required.

## License

This is a demonstration project. Use freely for learning and reference.

## Future Enhancements (If Extended)

While out of scope for this implementation, potential extensions include:

- B-tree indexes for range queries
- Query execution plans with EXPLAIN
- Additional data types (FLOAT, DATE, BLOB)
- LEFT/RIGHT JOIN support
- Basic aggregation functions (COUNT, SUM, AVG)
- Simple transaction support with rollback
- CSV import/export
- Disk persistence layer

## Design Philosophy

godb prioritizes:

1. **Readability** over performance
2. **Explicitness** over cleverness
3. **Simplicity** over feature completeness
4. **Testability** over convenience
5. **Separation of concerns** over monolithic design

The goal is to demonstrate solid engineering practices in a minimal database implementation.
