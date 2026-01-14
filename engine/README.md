# Engine Package

The `engine` package is the core of the `godb` database. It provides all the necessary components for creating and managing an in-memory database.

## Import

To use the `engine` package, import it as follows:

```go
import "godb/engine"
```

## Core Components

### Database

The `Database` struct is the main entry point for all database operations. It manages a collection of tables and provides methods for creating, dropping, and getting tables.

**Usage:**

```go
// Create a new database
db := engine.NewDatabase()

// Create a new table
schema := []engine.Column{
    {Name: "id", Type: engine.TypeInt, PrimaryKey: true},
    {Name: "name", Type: engine.TypeString, NotNull: true},
}
err := db.CreateTable("users", schema)
if err != nil {
    // Handle error
}

// Get a table
table, err := db.GetTable("users")
if err != nil {
    // Handle error
}
```

### CRUD Operations

The `Database` struct provides methods for performing CRUD (Create, Read, Update, Delete) operations on tables.

**Usage:**

```go
// Insert a new row
row := engine.Row{"id": 1, "name": "moses"}
err := db.Insert("users", row)
if err != nil {
    // Handle error
}

// Select rows
rows, err := db.Select("users", []string{"id", "name"}, nil)
if err != nil {
    // Handle error
}

// Update a row
updates := engine.Row{"name": "Alicia"}
condition := &engine.Condition{Column: "id", Operator: "=", Value: 1}
_, err = db.Update("users", updates, condition)
if err != nil {
    // Handle error
}

// Delete a row
condition = &engine.Condition{Column: "id", Operator: "=", Value: 1}
_, err = db.Delete("users", condition)
if err != nil {
    // Handle error
}
```

### Table, Row, Column, and Index

The `Table`, `Row`, `Column`, and `Index` structs are the building blocks of the database.

-   `Table`: Represents a table in the database, with a name, schema, and rows.
-   `Row`: Represents a single row in a table, as a map of column names to values.
-   `Column`: Represents a column in a table, with a name, type, and constraints.
-   `Index`: Represents an index on a column, for fast lookups.

## Errors

The `engine` package defines a set of custom error types to provide detailed information about database errors. These include `ErrTableNotFound`, `ErrPrimaryKeyViolation`, `ErrUniqueViolation`, and more.