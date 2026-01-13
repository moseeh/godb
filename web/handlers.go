package web

import (
	"encoding/json"
	"godb/engine"
	"net/http"
)

// Handler contains the database instance and HTTP handlers
type Handler struct {
	db *engine.Database
}

// NewHandler creates a new handler with a database instance
func NewHandler(db *engine.Database) *Handler {
	return &Handler{db: db}
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	row := engine.Row{
		"id":    req.ID,
		"name":  req.Name,
		"email": req.Email,
	}

	if err := h.db.Insert("users", row); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondSuccess(w, "User created successfully", 1)
}

// GetUsers handles GET /users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Select("users", nil, nil)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users := make([]UserResponse, 0, len(rows))
	for _, row := range rows {
		user := UserResponse{
			ID:    getInt(row, "id"),
			Name:  getString(row, "name"),
			Email: getString(row, "email"),
		}
		users = append(users, user)
	}

	respondJSON(w, users)
}

// CreatePost handles POST /posts
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	row := engine.Row{
		"id":      req.ID,
		"user_id": req.UserID,
		"title":   req.Title,
		"body":    req.Body,
	}

	if err := h.db.Insert("posts", row); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondSuccess(w, "Post created successfully", 1)
}

// GetPosts handles GET /posts (with JOIN to users)
func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Perform INNER JOIN between posts and users
	joinCondition := engine.JoinCondition{
		LeftColumn:  "user_id",
		RightColumn: "id",
	}

	rows, err := h.db.InnerJoin("posts", "users", joinCondition, nil)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts := make([]PostWithUserResponse, 0, len(rows))
	for _, row := range rows {
		post := PostWithUserResponse{
			PostID:    getInt(row, "posts.id"),
			PostTitle: getString(row, "posts.title"),
			PostBody:  getString(row, "posts.body"),
			UserID:    getInt(row, "users.id"),
			UserName:  getString(row, "users.name"),
			UserEmail: getString(row, "users.email"),
		}
		posts = append(posts, post)
	}

	respondJSON(w, posts)
}

// Helper functions

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func respondSuccess(w http.ResponseWriter, message string, count int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Count:   count,
	})
}

func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func getInt(row engine.Row, key string) int {
	if val, ok := row[key]; ok {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0
}

func getString(row engine.Row, key string) string {
	if val, ok := row[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
