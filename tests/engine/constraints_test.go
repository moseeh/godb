package engine_test

import (
	"godb/engine"
	"testing"
)

func TestPrimaryKeyConstraint(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}

	db.CreateTable("users", schema)

	// Insert first row
	row1 := engine.Row{"id": 1, "name": "Alice"}
	err := db.Insert("users", row1)
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	// Try to insert duplicate primary key
	row2 := engine.Row{"id": 1, "name": "Bob"}
	err = db.Insert("users", row2)

	if err == nil {
		t.Fatal("Expected primary key violation, got nil")
	}

	// Verify only one row exists
	rows, _ := db.Select("users", nil, nil)
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
}

func TestUniqueConstraint(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
		{Name: "name", Type: engine.TypeString},
	}

	db.CreateTable("users", schema)

	// Insert first row
	row1 := engine.Row{"id": 1, "email": "alice@example.com", "name": "Alice"}
	err := db.Insert("users", row1)
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	// Try to insert duplicate unique value
	row2 := engine.Row{"id": 2, "email": "alice@example.com", "name": "Alice2"}
	err = db.Insert("users", row2)

	if err == nil {
		t.Fatal("Expected unique constraint violation, got nil")
	}

	// Verify only one row exists
	rows, _ := db.Select("users", nil, nil)
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
}

func TestNotNullConstraint(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString, NotNull: true},
	}

	db.CreateTable("users", schema)

	// Try to insert row without required field
	row := engine.Row{"id": 1}
	err := db.Insert("users", row)

	if err == nil {
		t.Fatal("Expected missing required column error, got nil")
	}

	// Verify no rows exist
	rows, _ := db.Select("users", nil, nil)
	if len(rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(rows))
	}
}

func TestUpdateConstraintViolation(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
	}

	db.CreateTable("users", schema)

	// Insert two rows
	db.Insert("users", engine.Row{"id": 1, "email": "alice@example.com"})
	db.Insert("users", engine.Row{"id": 2, "email": "bob@example.com"})

	// Try to update second row with duplicate email
	updates := engine.Row{"email": "alice@example.com"}
	condition := &engine.Condition{
		Column:   "id",
		Operator: "=",
		Value:    2,
	}

	_, err := db.Update("users", updates, condition)

	if err == nil {
		t.Fatal("Expected unique constraint violation on update, got nil")
	}
}

func TestMultipleConstraints(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "username", Type: engine.TypeString, Unique: true, NotNull: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
	}

	db.CreateTable("users", schema)

	// Insert valid row
	row1 := engine.Row{"id": 1, "username": "alice", "email": "alice@example.com"}
	err := db.Insert("users", row1)
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	// Try to insert duplicate username
	row2 := engine.Row{"id": 2, "username": "alice", "email": "bob@example.com"}
	err = db.Insert("users", row2)
	if err == nil {
		t.Error("Expected unique constraint violation for username")
	}

	// Try to insert duplicate email
	row3 := engine.Row{"id": 3, "username": "bob", "email": "alice@example.com"}
	err = db.Insert("users", row3)
	if err == nil {
		t.Error("Expected unique constraint violation for email")
	}

	// Verify only one row exists
	rows, _ := db.Select("users", nil, nil)
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
}
