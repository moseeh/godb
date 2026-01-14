package engine_test

import (
	"godb/engine"
	"testing"
)

func TestSelectAll(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}

	db.CreateTable("users", schema)

	// Insert test data
	rows := []engine.Row{
		{"id": 1, "name": "moses"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Charlie"},
	}

	for _, row := range rows {
		db.Insert("users", row)
	}

	// Select all
	results, err := db.Select("users", nil, nil)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(results))
	}
}

func TestSelectWithCondition(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "age", Type: engine.TypeInt},
	}

	db.CreateTable("users", schema)

	// Insert test data
	rows := []engine.Row{
		{"id": 1, "age": 25},
		{"id": 2, "age": 30},
		{"id": 3, "age": 35},
	}

	for _, row := range rows {
		db.Insert("users", row)
	}

	// Select with condition
	condition := &engine.Condition{
		Column:   "age",
		Operator: ">",
		Value:    30,
	}

	results, err := db.Select("users", nil, condition)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 row, got %d", len(results))
	}

	if results[0]["age"] != 35 {
		t.Errorf("Expected age 35, got %v", results[0]["age"])
	}
}

func TestSelectWithIndex(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
	}

	db.CreateTable("users", schema)

	// Insert test data
	rows := []engine.Row{
		{"id": 1, "email": "moses@example.com"},
		{"id": 2, "email": "bob@example.com"},
		{"id": 3, "email": "charlie@example.com"},
	}

	for _, row := range rows {
		db.Insert("users", row)
	}

	// Select using indexed column (email)
	condition := &engine.Condition{
		Column:   "email",
		Operator: "=",
		Value:    "bob@example.com",
	}

	results, err := db.Select("users", nil, condition)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 row, got %d", len(results))
	}

	if results[0]["id"] != 2 {
		t.Errorf("Expected id 2, got %v", results[0]["id"])
	}
}

func TestSelectSpecificColumns(t *testing.T) {
	db := engine.NewDatabase()

	schema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
		{Name: "email", Type: engine.TypeString},
	}

	db.CreateTable("users", schema)

	row := engine.Row{"id": 1, "name": "moses", "email": "moses@example.com"}
	db.Insert("users", row)

	// Select only specific columns
	results, err := db.Select("users", []string{"name"}, nil)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 row, got %d", len(results))
	}

	// Should only have 'name' column
	if len(results[0]) != 1 {
		t.Errorf("Expected 1 column, got %d", len(results[0]))
	}

	if _, ok := results[0]["name"]; !ok {
		t.Error("Expected 'name' column to be present")
	}

	if _, ok := results[0]["email"]; ok {
		t.Error("Did not expect 'email' column to be present")
	}
}
