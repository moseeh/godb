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
		h.renderRowsWithTable(w, rows, c.TableName)

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

	h.renderRowsWithTable(w, rows, tableName)
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

// ColumnInfo contains column metadata for the results template
type ColumnInfo struct {
	Name         string
	IsPrimaryKey bool
}

func (h *Handler) renderRows(w http.ResponseWriter, rows []engine.Row) {
	h.renderRowsWithTable(w, rows, "")
}

func (h *Handler) renderRowsWithTable(w http.ResponseWriter, rows []engine.Row, tableName string) {
	if len(rows) == 0 {
		data := map[string]interface{}{
			"Success": true,
			"Message": "Query executed successfully - no rows returned",
			"Rows":    []engine.Row{},
		}
		h.renderResults(w, data, "")
		return
	}

	// Extract column names from first row
	var columnNames []string
	for col := range rows[0] {
		columnNames = append(columnNames, col)
	}

	// Build column info with primary key flags
	columns := make([]ColumnInfo, 0, len(columnNames))
	var pkColumn string

	// Try to get primary key info from the table
	if tableName != "" {
		if table, err := h.db.GetTable(tableName); err == nil {
			pkColumn = table.PrimaryKey()
		}
	}

	for _, colName := range columnNames {
		columns = append(columns, ColumnInfo{
			Name:         colName,
			IsPrimaryKey: colName == pkColumn,
		})
	}

	data := map[string]interface{}{
		"Rows":    rows,
		"Columns": columns,
	}
	h.renderResults(w, data, "")
}

// === UPDATE & DELETE HANDLERS ===

// UpdateTab renders the update data tab
func (h *Handler) UpdateTab(w http.ResponseWriter, r *http.Request) {
	tables := h.db.ListTables()
	data := map[string]interface{}{
		"Tables": tables,
	}
	if err := h.templates.ExecuteTemplate(w, "update", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteTab renders the delete data tab
func (h *Handler) DeleteTab(w http.ResponseWriter, r *http.Request) {
	tables := h.db.ListTables()
	data := map[string]interface{}{
		"Tables": tables,
	}
	if err := h.templates.ExecuteTemplate(w, "delete", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TableSchemaUpdate returns the schema for the update finder form
func (h *Handler) TableSchemaUpdate(w http.ResponseWriter, r *http.Request) {
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

	if err := h.templates.ExecuteTemplate(w, "update-finder", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TableSchemaDelete returns the schema for the delete condition form
func (h *Handler) TableSchemaDelete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.templates.ExecuteTemplate(w, "delete-condition", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FetchRow fetches a row for editing based on a condition
func (h *Handler) FetchRow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderUpdateEditor(w, nil, "Failed to parse form")
		return
	}

	tableName := r.FormValue("table_name")
	whereColumn := r.FormValue("where_column")
	whereValue := r.FormValue("where_value")

	table, err := h.db.GetTable(tableName)
	if err != nil {
		h.renderUpdateEditor(w, nil, err.Error())
		return
	}

	// Build condition with type conversion
	condition := h.buildCondition(table, whereColumn, "=", whereValue)

	// Fetch matching rows
	rows, err := h.db.Select(tableName, nil, condition)
	if err != nil {
		h.renderUpdateEditor(w, nil, err.Error())
		return
	}

	data := map[string]interface{}{
		"TableName":   tableName,
		"Schema":      table.Schema(),
		"WhereColumn": whereColumn,
		"WhereValue":  whereValue,
	}

	if len(rows) == 0 {
		data["NoRows"] = true
	} else if len(rows) > 1 {
		data["MultipleRows"] = true
		data["RowCount"] = len(rows)
	} else {
		data["RowData"] = rows[0]
	}

	h.renderUpdateEditor(w, data, "")
}

// PreviewDelete shows rows that will be deleted
func (h *Handler) PreviewDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderDeletePreview(w, nil, "Failed to parse form")
		return
	}

	tableName := r.FormValue("table_name")
	whereColumn := r.FormValue("where_column")
	whereOperator := r.FormValue("where_operator")
	whereValue := r.FormValue("where_value")

	table, err := h.db.GetTable(tableName)
	if err != nil {
		h.renderDeletePreview(w, nil, err.Error())
		return
	}

	// Build condition with type conversion
	condition := h.buildCondition(table, whereColumn, whereOperator, whereValue)

	// Fetch matching rows
	rows, err := h.db.Select(tableName, nil, condition)
	if err != nil {
		h.renderDeletePreview(w, nil, err.Error())
		return
	}

	// Build column info for display
	schema := table.Schema()
	columns := make([]map[string]interface{}, len(schema))
	for i, col := range schema {
		columns[i] = map[string]interface{}{
			"Name":       col.Name,
			"PrimaryKey": col.PrimaryKey,
		}
	}

	data := map[string]interface{}{
		"TableName":     tableName,
		"WhereColumn":   whereColumn,
		"WhereOperator": whereOperator,
		"WhereValue":    formatWhereValue(whereValue),
		"Columns":       columns,
		"Rows":          rows,
		"RowCount":      len(rows),
	}

	if len(rows) == 0 {
		data["NoRows"] = true
	}

	h.renderDeletePreview(w, data, "")
}

// BuildUpdate builds and executes an UPDATE statement from form data
func (h *Handler) BuildUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	tableName := r.FormValue("table_name")
	whereColumn := r.FormValue("where_column")
	whereValue := r.FormValue("where_value")

	table, err := h.db.GetTable(tableName)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	// Build updates from form values (excluding primary key)
	updates := make(engine.Row)
	pkColumn := table.PrimaryKey()

	for _, col := range table.Schema() {
		if col.Name == pkColumn {
			continue // Skip primary key
		}

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
			updates[col.Name] = intVal
		case engine.TypeBool:
			updates[col.Name] = value == "true"
		default:
			updates[col.Name] = value
		}
	}

	// Build condition
	condition := h.buildCondition(table, whereColumn, "=", whereValue)

	// Execute UPDATE
	rowsAffected, err := h.db.Update(tableName, updates, condition)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	h.renderSuccess(w, fmt.Sprintf("%d row(s) updated successfully", rowsAffected))
}

// BuildDelete builds and executes a DELETE statement from form data
func (h *Handler) BuildDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	tableName := r.FormValue("table_name")
	whereColumn := r.FormValue("where_column")
	whereOperator := r.FormValue("where_operator")
	whereValue := r.FormValue("where_value")

	table, err := h.db.GetTable(tableName)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	// Build condition
	condition := h.buildCondition(table, whereColumn, whereOperator, whereValue)

	// Execute DELETE
	rowsAffected, err := h.db.Delete(tableName, condition)
	if err != nil {
		h.renderResults(w, nil, err.Error())
		return
	}

	h.renderSuccess(w, fmt.Sprintf("%d row(s) deleted successfully", rowsAffected))
}

// Helper: build condition with proper type conversion
func (h *Handler) buildCondition(table *engine.Table, column, operator, value string) *engine.Condition {
	// Find column type
	var colType engine.ColumnType
	for _, col := range table.Schema() {
		if col.Name == column {
			colType = col.Type
			break
		}
	}

	var typedValue interface{}
	switch colType {
	case engine.TypeInt:
		if intVal, err := strconv.Atoi(value); err == nil {
			typedValue = intVal
		} else {
			typedValue = value
		}
	case engine.TypeBool:
		typedValue = value == "true"
	default:
		typedValue = value
	}

	return &engine.Condition{
		Column:   column,
		Operator: operator,
		Value:    typedValue,
	}
}

// Helper: format where value for SQL display
func formatWhereValue(value string) string {
	// If it looks like a number, return as-is
	if _, err := strconv.Atoi(value); err == nil {
		return value
	}
	// If it's a boolean, return as-is
	if value == "true" || value == "false" {
		return value
	}
	// Otherwise, wrap in quotes
	return "'" + value + "'"
}

// Helper: render update editor template
func (h *Handler) renderUpdateEditor(w http.ResponseWriter, data map[string]interface{}, errorMsg string) {
	if data == nil {
		data = make(map[string]interface{})
	}
	if errorMsg != "" {
		data["Error"] = errorMsg
	}
	if err := h.templates.ExecuteTemplate(w, "update-editor", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper: render delete preview template
func (h *Handler) renderDeletePreview(w http.ResponseWriter, data map[string]interface{}, errorMsg string) {
	if data == nil {
		data = make(map[string]interface{})
	}
	if errorMsg != "" {
		data["Error"] = errorMsg
	}
	if err := h.templates.ExecuteTemplate(w, "delete-preview", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

