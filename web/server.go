package web

import (
	"fmt"
	"godb/engine"
	"log"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	db   *engine.Database
	addr string
}

// NewServer creates a new HTTP server
func NewServer(addr string) *Server {
	return &Server{
		db:   engine.NewDatabase(),
		addr: addr,
	}
}

// Initialize sets up the database schema for the demo
func (s *Server) Initialize() error {
	// Create users table
	usersSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "name", Type: engine.TypeString, NotNull: true},
		{Name: "email", Type: engine.TypeString, Unique: true},
	}

	if err := s.db.CreateTable("users", usersSchema); err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Create posts table
	postsSchema := []engine.Column{
		{Name: "id", Type: engine.TypeInt, PrimaryKey: true},
		{Name: "user_id", Type: engine.TypeInt, NotNull: true},
		{Name: "title", Type: engine.TypeString, NotNull: true},
		{Name: "body", Type: engine.TypeString},
	}

	if err := s.db.CreateTable("posts", postsSchema); err != nil {
		return fmt.Errorf("failed to create posts table: %v", err)
	}

	// Create index on posts.user_id for efficient joins
	table, _ := s.db.GetTable("posts")
	if err := table.CreateIndex("user_id"); err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	handler := NewHandler(s.db)

	// Register routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateUser(w, r)
		case http.MethodGet:
			handler.GetUsers(w, r)
		default:
			respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreatePost(w, r)
		case http.MethodGet:
			handler.GetPosts(w, r)
		default:
			respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Starting godb web server on %s", s.addr)
	log.Println("Available endpoints:")
	log.Println("  POST /users - Create a user")
	log.Println("  GET  /users - List all users")
	log.Println("  POST /posts - Create a post")
	log.Println("  GET  /posts - List all posts (with JOIN to users)")

	return http.ListenAndServe(s.addr, nil)
}
