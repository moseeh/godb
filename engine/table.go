package engine

// Table represents a database table with schema, data, and indexes
type Table struct {
	name       string
	schema     []Column
	rows       []Row
	primaryKey string
	indexes    map[string]*Index // column name -> index
}

// NewTable creates a new table with the given schema
func NewTable(name string, schema []Column) *Table {
	table := &Table{
		name:    name,
		schema:  schema,
		rows:    make([]Row, 0),
		indexes: make(map[string]*Index),
	}

	// Identify primary key and create indexes
	for _, col := range schema {
		if col.PrimaryKey {
			table.primaryKey = col.Name
			table.CreateIndex(col.Name)
		} else if col.Unique {
			table.CreateIndex(col.Name)
		}
	}

	return table
}

// Name returns the table name
func (t *Table) Name() string {
	return t.name
}

// Schema returns the table schema
func (t *Table) Schema() []Column {
	return t.schema
}

// Rows returns all rows in the table
func (t *Table) Rows() []Row {
	return t.rows
}

// PrimaryKey returns the primary key column name
func (t *Table) PrimaryKey() string {
	return t.primaryKey
}

// CreateIndex creates an index on a column
func (t *Table) CreateIndex(columnName string) error {
	// Check if column exists
	if !t.hasColumn(columnName) {
		return ErrColumnNotFound{
			TableName:  t.name,
			ColumnName: columnName,
		}
	}

	// Check if index already exists
	if _, exists := t.indexes[columnName]; exists {
		return nil // Index already exists
	}

	// Create index
	idx := NewIndex(columnName)

	// Build index from existing rows
	for rowIdx, row := range t.rows {
		if value, ok := row.Get(columnName); ok {
			idx.Add(value, rowIdx)
		}
	}

	t.indexes[columnName] = idx
	return nil
}

// GetIndex returns the index for a column if it exists
func (t *Table) GetIndex(columnName string) (*Index, bool) {
	idx, ok := t.indexes[columnName]
	return idx, ok
}

// hasColumn checks if a column exists in the table schema
func (t *Table) hasColumn(columnName string) bool {
	for _, col := range t.schema {
		if col.Name == columnName {
			return true
		}
	}
	return false
}

// hasPrimaryKeyValue checks if a primary key value already exists
func (t *Table) hasPrimaryKeyValue(value interface{}) bool {
	if t.primaryKey == "" {
		return false
	}

	idx, hasIndex := t.indexes[t.primaryKey]
	if hasIndex {
		return idx.Has(value)
	}

	// Fallback: linear scan
	for _, row := range t.rows {
		if rowValue, ok := row.Get(t.primaryKey); ok && rowValue == value {
			return true
		}
	}
	return false
}

// hasUniqueValue checks if a unique column value already exists
func (t *Table) hasUniqueValue(columnName string, value interface{}) bool {
	idx, hasIndex := t.indexes[columnName]
	if hasIndex {
		return idx.Has(value)
	}

	// Fallback: linear scan
	for _, row := range t.rows {
		if rowValue, ok := row.Get(columnName); ok && rowValue == value {
			return true
		}
	}
	return false
}

// addRow adds a row to the table and updates indexes
func (t *Table) addRow(row Row) int {
	rowIndex := len(t.rows)
	t.rows = append(t.rows, row)

	// Update all indexes
	for colName, idx := range t.indexes {
		if value, ok := row.Get(colName); ok {
			idx.Add(value, rowIndex)
		}
	}

	return rowIndex
}

// updateRow updates a row at a given index and updates indexes
func (t *Table) updateRow(rowIndex int, newRow Row) {
	oldRow := t.rows[rowIndex]

	// Update indexes
	for colName, idx := range t.indexes {
		oldValue, _ := oldRow.Get(colName)
		newValue, _ := newRow.Get(colName)
		if oldValue != newValue {
			idx.Update(oldValue, newValue, rowIndex)
		}
	}

	t.rows[rowIndex] = newRow
}

// deleteRow removes a row at a given index and updates indexes
func (t *Table) deleteRow(rowIndex int) {
	row := t.rows[rowIndex]

	// Update indexes
	for colName, idx := range t.indexes {
		if value, ok := row.Get(colName); ok {
			idx.Remove(value, rowIndex)
		}
	}

	// Remove row by replacing with last element and truncating
	lastIndex := len(t.rows) - 1
	if rowIndex != lastIndex {
		t.rows[rowIndex] = t.rows[lastIndex]

		// Update indexes for moved row
		for colName, idx := range t.indexes {
			if value, ok := t.rows[rowIndex].Get(colName); ok {
				idx.Remove(value, lastIndex)
				idx.Add(value, rowIndex)
			}
		}
	}

	t.rows = t.rows[:lastIndex]
}
