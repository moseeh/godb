package engine

// Condition represents a WHERE clause condition
type Condition struct {
	Column   string
	Operator string // "=", "!=", ">", "<", ">=", "<="
	Value    interface{}
}

// Insert adds a new row to a table
func (db *Database) Insert(tableName string, row Row) error {
	table, err := db.GetTable(tableName)
	if err != nil {
		return err
	}

	// Validate constraints
	checker := NewConstraintChecker(table)
	if err := checker.ValidateInsert(row); err != nil {
		return err
	}

	// Add row to table
	table.addRow(row)
	return nil
}

// Select retrieves rows from a table with optional filtering
func (db *Database) Select(tableName string, columns []string, condition *Condition) ([]Row, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	// Get candidate rows
	var candidateIndices []int
	useIndex := false

	// Try to use index if condition is on an indexed column with equality
	if condition != nil && condition.Operator == "=" {
		if idx, hasIdx := table.GetIndex(condition.Column); hasIdx {
			candidateIndices = idx.Lookup(condition.Value)
			useIndex = true
		}
	}

	// If no index used, scan all rows
	if !useIndex {
		candidateIndices = make([]int, len(table.rows))
		for i := range table.rows {
			candidateIndices[i] = i
		}
	}

	// Filter rows based on condition
	var results []Row
	for _, idx := range candidateIndices {
		if idx >= len(table.rows) {
			continue // Skip invalid indices
		}
		row := table.rows[idx]

		// Apply condition if present
		if condition != nil {
			if !evaluateCondition(row, condition) {
				continue
			}
		}

		// Project columns
		resultRow := projectRow(row, columns, table.schema)
		results = append(results, resultRow)
	}

	return results, nil
}

// Update modifies rows in a table that match the condition
func (db *Database) Update(tableName string, updates Row, condition *Condition) (int, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return 0, err
	}

	checker := NewConstraintChecker(table)
	rowsAffected := 0

	// Find rows to update
	for i := 0; i < len(table.rows); i++ {
		row := table.rows[i]

		// Check if row matches condition
		if condition != nil && !evaluateCondition(row, condition) {
			continue
		}

		// Create updated row
		newRow := row.Copy()
		for col, val := range updates {
			newRow.Set(col, val)
		}

		// Validate constraints
		if err := checker.ValidateUpdate(row, newRow); err != nil {
			return rowsAffected, err
		}

		// Update the row
		table.updateRow(i, newRow)
		rowsAffected++
	}

	return rowsAffected, nil
}

// Delete removes rows from a table that match the condition
func (db *Database) Delete(tableName string, condition *Condition) (int, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return 0, err
	}

	rowsAffected := 0

	// Iterate backwards to avoid index issues when deleting
	for i := len(table.rows) - 1; i >= 0; i-- {
		row := table.rows[i]

		// Check if row matches condition
		if condition != nil && !evaluateCondition(row, condition) {
			continue
		}

		// Delete the row
		table.deleteRow(i)
		rowsAffected++
	}

	return rowsAffected, nil
}

// evaluateCondition checks if a row satisfies a condition
func evaluateCondition(row Row, cond *Condition) bool {
	value, ok := row.Get(cond.Column)
	if !ok {
		return false
	}

	switch cond.Operator {
	case "=":
		return value == cond.Value
	case "!=":
		return value != cond.Value
	case ">":
		return compareValues(value, cond.Value) > 0
	case "<":
		return compareValues(value, cond.Value) < 0
	case ">=":
		return compareValues(value, cond.Value) >= 0
	case "<=":
		return compareValues(value, cond.Value) <= 0
	default:
		return false
	}
}

// compareValues compares two values for ordering
func compareValues(a, b interface{}) int {
	switch av := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		}
	case string:
		if bv, ok := b.(string); ok {
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		}
	}
	return 0
}

// projectRow extracts specified columns from a row
// If columns is empty, returns all columns
func projectRow(row Row, columns []string, schema []Column) Row {
	if len(columns) == 0 {
		return row.Copy()
	}

	result := make(Row)
	for _, col := range columns {
		if value, ok := row.Get(col); ok {
			result.Set(col, value)
		}
	}
	return result
}
