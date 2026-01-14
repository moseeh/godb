package engine_test

import (
	"godb/engine"
	"testing"
)

func TestInsertBasic(t *testing.T) {
	db := engine.NewDatabase()

	// Create table
	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString, NotNull: true},
	}

	err := db.CreateTable("users", schema)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert row
	row := engine.Row{
		"id":   1,
		"name": "moses",
	}

	err = db.Insert("users", row)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}

	// Verify insertion
	rows, err := db.Select("users", nil, nil)
	if err != nil {
		t.Fatalf("Failed to select: %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
}

func TestInsertDuplicatePrimaryKey(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}

	db.CreateTable("users", schema)

	// Insert first row
	row1 := engine.Row{"id": 1, "name": "moses"}
	err := db.Insert("users", row1)
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	// Insert duplicate primary key
	row2 := engine.Row{"id": 1, "name": "Bob"}
	err = db.Insert("users", row2)

	// Should get primary key violation error
	if err == nil {
		t.Error("Expected primary key violation error, got nil")
	}

	if _, ok := err.(engine.ErrPrimaryKeyViolation); !ok {
		t.Errorf("Expected ErrPrimaryKeyViolation, got %T", err)
	}
}

func TestInsertUniqueViolation(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
	}

	db.CreateTable("users", schema)

	// Insert first row
	row1 := engine.Row{"id": 1, "email": "moses@example.com"}
	err := db.Insert("users", row1)
	if err != nil {
		t.Fatalf("First insert failed: %v", err)
	}

	// Insert duplicate unique value
	row2 := engine.Row{"id": 2, "email": "moses@example.com"}
	err = db.Insert("users", row2)

	// Should get unique violation error
	if err == nil {
		t.Error("Expected unique violation error, got nil")
	}

	if _, ok := err.(engine.ErrUniqueViolation); !ok {
		t.Errorf("Expected ErrUniqueViolation, got %T", err)
	}
}

func TestInsertMissingRequiredColumn(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString, NotNull: true},
	}

	db.CreateTable("users", schema)

	// Insert row without required column
	row := engine.Row{"id": 1}
	err := db.Insert("users", row)

	// Should get missing column error
	if err == nil {
		t.Error("Expected missing required column error, got nil")
	}

	if _, ok := err.(engine.ErrMissingRequiredColumn); !ok {
		t.Errorf("Expected ErrMissingRequiredColumn, got %T", err)
	}
}
