package parser_test

import (
	"godb/parser"
	"testing"
)

func TestParseCreateTable(t *testing.T) {
	input := "CREATE TABLE users (id INT PRIMARY KEY, name STRING NOT NULL, email STRING UNIQUE)"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	createCmd, ok := cmd.(*parser.CreateTableCommand)
	if !ok {
		t.Fatalf("Expected CreateTableCommand, got %T", cmd)
	}

	if createCmd.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", createCmd.TableName)
	}

	if len(createCmd.Columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(createCmd.Columns))
	}

	// Check first column
	if createCmd.Columns[0].Name != "id" {
		t.Errorf("Expected column name 'id', got '%s'", createCmd.Columns[0].Name)
	}

	if !createCmd.Columns[0].PrimaryKey {
		t.Error("Expected 'id' to be primary key")
	}

	// Check second column
	if !createCmd.Columns[1].NotNull {
		t.Error("Expected 'name' to be NOT NULL")
	}

	// Check third column
	if !createCmd.Columns[2].Unique {
		t.Error("Expected 'email' to be UNIQUE")
	}
}

func TestParseInsert(t *testing.T) {
	input := "INSERT INTO users (id, name) VALUES (1, 'moses')"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	insertCmd, ok := cmd.(*parser.InsertCommand)
	if !ok {
		t.Fatalf("Expected InsertCommand, got %T", cmd)
	}

	if insertCmd.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", insertCmd.TableName)
	}

	if len(insertCmd.Values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(insertCmd.Values))
	}

	if insertCmd.Values["id"] != 1 {
		t.Errorf("Expected id=1, got %v", insertCmd.Values["id"])
	}

	if insertCmd.Values["name"] != "moses" {
		t.Errorf("Expected name='moses', got %v", insertCmd.Values["name"])
	}
}

func TestParseSelect(t *testing.T) {
	input := "SELECT name, email FROM users WHERE id = 1"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectCmd, ok := cmd.(*parser.SelectCommand)
	if !ok {
		t.Fatalf("Expected SelectCommand, got %T", cmd)
	}

	if selectCmd.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", selectCmd.TableName)
	}

	if len(selectCmd.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(selectCmd.Columns))
	}

	if selectCmd.Condition == nil {
		t.Fatal("Expected condition to be present")
	}

	if selectCmd.Condition.Column != "id" {
		t.Errorf("Expected condition column 'id', got '%s'", selectCmd.Condition.Column)
	}

	if selectCmd.Condition.Operator != "=" {
		t.Errorf("Expected operator '=', got '%s'", selectCmd.Condition.Operator)
	}

	if selectCmd.Condition.Value != 1 {
		t.Errorf("Expected value 1, got %v", selectCmd.Condition.Value)
	}
}

func TestParseSelectAll(t *testing.T) {
	input := "SELECT * FROM users"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	selectCmd, ok := cmd.(*parser.SelectCommand)
	if !ok {
		t.Fatalf("Expected SelectCommand, got %T", cmd)
	}

	// nil columns means select all
	if selectCmd.Columns != nil {
		t.Errorf("Expected nil columns for SELECT *, got %v", selectCmd.Columns)
	}
}

func TestParseUpdate(t *testing.T) {
	input := "UPDATE users SET name = 'Bob', email = 'bob@example.com' WHERE id = 1"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	updateCmd, ok := cmd.(*parser.UpdateCommand)
	if !ok {
		t.Fatalf("Expected UpdateCommand, got %T", cmd)
	}

	if updateCmd.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", updateCmd.TableName)
	}

	if len(updateCmd.Updates) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(updateCmd.Updates))
	}

	if updateCmd.Updates["name"] != "Bob" {
		t.Errorf("Expected name='Bob', got %v", updateCmd.Updates["name"])
	}

	if updateCmd.Condition == nil {
		t.Fatal("Expected condition to be present")
	}
}

func TestParseDelete(t *testing.T) {
	input := "DELETE FROM users WHERE id = 1"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	deleteCmd, ok := cmd.(*parser.DeleteCommand)
	if !ok {
		t.Fatalf("Expected DeleteCommand, got %T", cmd)
	}

	if deleteCmd.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", deleteCmd.TableName)
	}

	if deleteCmd.Condition == nil {
		t.Fatal("Expected condition to be present")
	}
}

func TestParseJoin(t *testing.T) {
	input := "SELECT * FROM posts INNER JOIN users ON posts.user_id = users.id"
	p := parser.NewParser(input)
	cmd, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	joinCmd, ok := cmd.(*parser.JoinCommand)
	if !ok {
		t.Fatalf("Expected JoinCommand, got %T", cmd)
	}

	if joinCmd.LeftTable != "posts" {
		t.Errorf("Expected left table 'posts', got '%s'", joinCmd.LeftTable)
	}

	if joinCmd.RightTable != "users" {
		t.Errorf("Expected right table 'users', got '%s'", joinCmd.RightTable)
	}

	if joinCmd.LeftColumn != "user_id" {
		t.Errorf("Expected left column 'user_id', got '%s'", joinCmd.LeftColumn)
	}

	if joinCmd.RightColumn != "id" {
		t.Errorf("Expected right column 'id', got '%s'", joinCmd.RightColumn)
	}
}
