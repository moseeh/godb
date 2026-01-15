package web

import (
	"fmt"
	"godb/engine"
	"html/template"
	"log"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	db        *engine.Database
	addr      string
	templates *template.Template
}

// NewServer creates a new HTTP server
func NewServer(addr string) *Server {
	// Parse all templates
	templates := template.Must(template.ParseGlob("web/templates/*.html"))

	return &Server{
		db:        engine.NewDatabase(),
		addr:      addr,
		templates: templates,
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
	handler := NewHandler(s.db, s.templates)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// UI routes
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/tabs/console", handler.ConsoleTab)
	http.HandleFunc("/tabs/create", handler.CreateTab)
	http.HandleFunc("/tabs/insert", handler.InsertTab)
	http.HandleFunc("/tabs/query", handler.QueryTab)
	http.HandleFunc("/tabs/update", handler.UpdateTab)
	http.HandleFunc("/tabs/delete", handler.DeleteTab)

	// Wizard routes
	http.HandleFunc("/wizard/create/step2", handler.CreateStep2)
	http.HandleFunc("/wizard/create/review", handler.CreateReview)

	// Action routes
	http.HandleFunc("/execute", handler.ExecuteSQL)
	http.HandleFunc("/table-schema", handler.TableSchema)
	http.HandleFunc("/build-insert", handler.BuildInsert)
	http.HandleFunc("/build-select", handler.BuildSelect)
	http.HandleFunc("/build-update", handler.BuildUpdate)
	http.HandleFunc("/build-delete", handler.BuildDelete)

	// Update/Delete helper routes
	http.HandleFunc("/table-schema-update", handler.TableSchemaUpdate)
	http.HandleFunc("/table-schema-delete", handler.TableSchemaDelete)
	http.HandleFunc("/fetch-row", handler.FetchRow)
	http.HandleFunc("/preview-delete", handler.PreviewDelete)

	// Legacy API routes (kept for backward compatibility)
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
	log.Println("Available interfaces:")
	log.Println("  Web UI:  http://localhost:8080/")
	log.Println("  API:     POST /users, GET /users, POST /posts, GET /posts")

	return http.ListenAndServe(s.addr, nil)
}
