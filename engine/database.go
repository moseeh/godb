package engine

import "sync"

// Database represents the in-memory database with multiple tables
type Database struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

// NewDatabase creates a new empty database
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table with the given schema
func (db *Database) CreateTable(name string, schema []Column) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; exists {
		return ErrTableAlreadyExists{TableName: name}
	}

	table := NewTable(name, schema)
	db.tables[name] = table
	return nil
}

// GetTable retrieves a table by name
func (db *Database) GetTable(name string) (*Table, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	table, exists := db.tables[name]
	if !exists {
		return nil, ErrTableNotFound{TableName: name}
	}

	return table, nil
}

// DropTable removes a table from the database
func (db *Database) DropTable(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; !exists {
		return ErrTableNotFound{TableName: name}
	}

	delete(db.tables, name)
	return nil
}

// ListTables returns the names of all tables in the database
func (db *Database) ListTables() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	names := make([]string, 0, len(db.tables))
	for name := range db.tables {
		names = append(names, name)
	}
	return names
}

// TableExists checks if a table exists in the database
func (db *Database) TableExists(name string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, exists := db.tables[name]
	return exists
}
