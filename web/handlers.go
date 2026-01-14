package web

import (
	"encoding/json"
	"fmt"
	"godb/engine"
	"godb/parser"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Handler contains the database instance and HTTP handlers
type Handler struct {
	db        *engine.Database
	templates *template.Template
}

// NewHandler creates a new handler with a database instance
func NewHandler(db *engine.Database, templates *template.Template) *Handler {
	return &Handler{
		db:        db,
		templates: templates,
	}
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

// === NEW UI HANDLERS ===

// Index renders the main page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "layout.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ConsoleTab renders the SQL console tab
func (h *Handler) ConsoleTab(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "console", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateTab renders the create table wizard tab
func (h *Handler) CreateTab(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	if err := h.templates.ExecuteTemplate(w, "create", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// InsertTab renders the insert data tab
func (h *Handler) InsertTab(w http.ResponseWriter, r *http.Request) {
	tables := h.db.ListTables()
	data := map[string]interface{}{
		"Tables": tables,
	}
	if err := h.templates.ExecuteTemplate(w, "insert", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// QueryTab renders the query data tab
func (h *Handler) QueryTab(w http.ResponseWriter, r *http.Request) {
	tables := h.db.ListTables()
	data := map[string]interface{}{
		"Tables": tables,
	}
	if err := h.templates.ExecuteTemplate(w, "query", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateStep2 handles step 2 of the create table wizard
func (h *Handler) CreateStep2(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tableName := r.FormValue("table_name")
	data := map[string]interface{}{
		"Step":      "step2",
		"TableName": tableName,
		"Columns":   []map[string]interface{}{},
	}

	if err := h.templates.ExecuteTemplate(w, "create-step2", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateReview handles the review step of create table wizard
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tableName := r.FormValue("table_name")

	// Parse columns from form
	var columns []string
	i := 0
	for {
		colName := r.FormValue(fmt.Sprintf("col_name_%d", i))
		if colName == "" {
			break
		}

		colType := r.FormValue(fmt.Sprintf("col_type_%d", i))
		isPK := r.FormValue(fmt.Sprintf("col_pk_%d", i)) == "on"
		isUnique := r.FormValue(fmt.Sprintf("col_unique_%d", i)) == "on"
		isNotNull := r.FormValue(fmt.Sprintf("col_notnull_%d", i)) == "on"

		colDef := colName + " " + colType
		if isPK {
			colDef += " PRIMARY KEY"
		}
		if isUnique {
			colDef += " UNIQUE"
		}
		if isNotNull {
			colDef += " NOT NULL"
		}

		columns = append(columns, colDef)
		i++
	}

	sql := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))

	data := map[string]interface{}{
		"Step":      "review",
		"TableName": tableName,
		"SQL":       sql,
	}

	if err := h.templates.ExecuteTemplate(w, "create-review", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ExecuteSQL executes a SQL command and returns results
func (h *Handler) ExecuteSQL(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	sql := r.FormValue("sql")
	if sql == "" {
		h.renderResults(w, nil, "SQL command is required")
		return
	}

	// Parse the SQL
	p := parser.NewParser(sql)
	cmd, err := p.Parse()
	if err != nil {
		h.renderResults(w, nil, fmt.Sprintf("Parse error: %v", err))
		return
	}

	// Execute based on command type
	switch c := cmd.(type) {
	case *parser.CreateTableCommand:
		err = h.db.CreateTable(c.TableName, c.Columns)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderSuccess(w, "Table created successfully")

	case *parser.InsertCommand:
		err = h.db.Insert(c.TableName, c.Values)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderSuccess(w, "Row inserted successfully")

	case *parser.SelectCommand:
		rows, err := h.db.Select(c.TableName, c.Columns, c.Condition)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderRows(w, rows)

	case *parser.UpdateCommand:
		rowsAffected, err := h.db.Update(c.TableName, c.Updates, c.Condition)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderSuccess(w, fmt.Sprintf("%d row(s) updated", rowsAffected))

	case *parser.DeleteCommand:
		rowsAffected, err := h.db.Delete(c.TableName, c.Condition)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderSuccess(w, fmt.Sprintf("%d row(s) deleted", rowsAffected))

	case *parser.JoinCommand:
		joinCondition := engine.JoinCondition{
			LeftColumn:  c.LeftColumn,
			RightColumn: c.RightColumn,
		}
		rows, err := h.db.InnerJoin(c.LeftTable, c.RightTable, joinCondition, c.SelectColumns)
		if err != nil {
			h.renderResults(w, nil, err.Error())
			return
		}
		h.renderRows(w, rows)

	default:
		h.renderResults(w, nil, "Unknown command type")
	}
}

// TableSchema returns the schema for a table (for dynamic form generation)
func (h *Handler) TableSchema(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	if tableName == "" {
		http.Error(w, "Table name required", http.StatusBadRequest)
		return
	}

	table, err := h.db.GetTable(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"TableName": tableName,
		"Schema":    table.Schema(),
	}

	if err := h.templates.ExecuteTemplate(w, "insert-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// BuildInsert builds and executes an INSERT statement from form data
func (h *Handler) BuildInsert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	tableName := r.FormValue("table_name")
	table, err := h.db.GetTable(tableName)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	// Build row from form values
	row := make(engine.Row)
	for _, col := range table.Schema() {
		value := r.FormValue(col.Name)
		if value == "" {
			continue
		}

		switch col.Type {
		case engine.TypeInt:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				h.renderResults(w, nil, fmt.Sprintf("Invalid integer for %s: %s", col.Name, value))
				return
			}
			row[col.Name] = intVal
		case engine.TypeBool:
			row[col.Name] = value == "true"
		default:
			row[col.Name] = value
		}
	}

	// Execute INSERT
	if err := h.db.Insert(tableName, row); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	h.renderSuccess(w, "Row inserted successfully")
}

// BuildSelect builds and executes a SELECT statement from form data
func (h *Handler) BuildSelect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	tableName := r.FormValue("table")
	columnsStr := r.FormValue("columns")
	useWhere := r.FormValue("use_where") == "on"

	if tableName == "" {
		h.renderResults(w, nil, "Table name is required")
		return
	}

	// Parse columns
	var columns []string
	if columnsStr != "" && columnsStr != "*" {
		for _, col := range strings.Split(columnsStr, ",") {
			columns = append(columns, strings.TrimSpace(col))
		}
	}

	// Build WHERE condition
	var condition *engine.Condition
	if useWhere {
		column := r.FormValue("where_column")
		operator := r.FormValue("where_operator")
		value := r.FormValue("where_value")

		if column != "" && value != "" {
			// Try to parse as int
			if intVal, err := strconv.Atoi(value); err == nil {
				condition = &engine.Condition{
					Column:   column,
					Operator: operator,
					Value:    intVal,
				}
			} else if value == "true" || value == "false" {
				condition = &engine.Condition{
					Column:   column,
					Operator: operator,
					Value:    value == "true",
				}
			} else {
				condition = &engine.Condition{
					Column:   column,
					Operator: operator,
					Value:    value,
				}
			}
		}
	}

	// Execute SELECT
	rows, err := h.db.Select(tableName, columns, condition)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	h.renderRows(w, rows)
}

// Helper functions for rendering results

func (h *Handler) renderResults(w http.ResponseWriter, data map[string]interface{}, errorMsg string) {
	if data == nil {
		data = make(map[string]interface{})
	}
	if errorMsg != "" {
		data["Error"] = errorMsg
	}
	if err := h.templates.ExecuteTemplate(w, "results", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) renderSuccess(w http.ResponseWriter, message string) {
	data := map[string]interface{}{
		"Success": true,
		"Message": message,
	}
	h.renderResults(w, data, "")
}

func (h *Handler) renderRows(w http.ResponseWriter, rows []engine.Row) {
	if len(rows) == 0 {
		data := map[string]interface{}{
			"Success": true,
			"Message": "Query executed successfully",
			"Rows":    []engine.Row{},
		}
		h.renderResults(w, data, "")
		return
	}

	// Extract column names from first row
	var columns []string
	for col := range rows[0] {
		columns = append(columns, col)
	}

	data := map[string]interface{}{
		"Rows":    rows,
		"Columns": columns,
	}
	h.renderResults(w, data, "")
}

