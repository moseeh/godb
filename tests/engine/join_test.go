package engine_test

import (
	"godb/engine"
	"testing"
)

func TestInnerJoinBasic(t *testing.T) {
	db := engine.NewDatabase()

	// Create users table
	usersSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}
	db.CreateTable("users", usersSchema)

	// Create posts table
	postsSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "user_id", Type: engine.TypeInt},
		{Name: "title", Type: engine.TypeString},
	}
	db.CreateTable("posts", postsSchema)

	// Insert users
	db.Insert("users", engine.Row{"id": 1, "name": "moses"})
	db.Insert("users", engine.Row{"id": 2, "name": "Bob"})

	// Insert posts
	db.Insert("posts", engine.Row{"id": 1, "user_id": 1, "title": "Post 1"})
	db.Insert("posts", engine.Row{"id": 2, "user_id": 1, "title": "Post 2"})
	db.Insert("posts", engine.Row{"id": 3, "user_id": 2, "title": "Post 3"})

	// Perform join
	joinCondition := engine.JoinCondition{
		LeftColumn:  "user_id",
		RightColumn: "id",
	}

	results, err := db.InnerJoin("posts", "users", joinCondition, nil)
	if err != nil {
		t.Fatalf("Join failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 joined rows, got %d", len(results))
	}

	// Verify first result
	if results[0]["posts.title"] == nil {
		t.Error("Expected posts.title to be present")
	}
	if results[0]["users.name"] == nil {
		t.Error("Expected users.name to be present")
	}
}

func TestInnerJoinWithIndex(t *testing.T) {
	db := engine.NewDatabase()

	// Create tables
	usersSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}
	db.CreateTable("users", usersSchema)

	postsSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "user_id", Type: engine.TypeInt},
		{Name: "title", Type: engine.TypeString},
	}
	db.CreateTable("posts", postsSchema)

	// Create index on user_id for efficient join
	table, _ := db.GetTable("users")
	table.CreateIndex("id")

	// Insert data
	db.Insert("users", engine.Row{"id": 1, "name": "moses"})
	db.Insert("posts", engine.Row{"id": 1, "user_id": 1, "title": "Post 1"})

	// Join should use index
	joinCondition := engine.JoinCondition{
		LeftColumn:  "user_id",
		RightColumn: "id",
	}

	results, err := db.InnerJoin("posts", "users", joinCondition, nil)
	if err != nil {
		t.Fatalf("Join failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 joined row, got %d", len(results))
	}
}

func TestInnerJoinNoMatches(t *testing.T) {
	db := engine.NewDatabase()

	// Create tables
	usersSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString},
	}
	db.CreateTable("users", usersSchema)

	postsSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "user_id", Type: engine.TypeInt},
		{Name: "title", Type: engine.TypeString},
	}
	db.CreateTable("posts", postsSchema)

	// Insert users
	db.Insert("users", engine.Row{"id": 1, "name": "moses"})

	// Insert posts with non-matching user_id
	db.Insert("posts", engine.Row{"id": 1, "user_id": 999, "title": "Post 1"})

	// Perform join
	joinCondition := engine.JoinCondition{
		LeftColumn:  "user_id",
		RightColumn: "id",
	}

	results, err := db.InnerJoin("posts", "users", joinCondition, nil)
	if err != nil {
		t.Fatalf("Join failed: %v", err)
	}

	// Should return no rows since there's no match
	if len(results) != 0 {
		t.Errorf("Expected 0 joined rows, got %d", len(results))
	}
}
