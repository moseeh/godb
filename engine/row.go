package engine

// Row represents a single row in a table
// Each row is a map of column name to value
type Row map[string]interface{}

// Copy creates a deep copy of the row
func (r Row) Copy() Row {
	copy := make(Row, len(r))
	for k, v := range r {
		copy[k] = v
	}
	return copy
}

// Get retrieves a value from the row by column name
func (r Row) Get(column string) (interface{}, bool) {
	val, ok := r[column]
	return val, ok
}

// Set sets a value in the row for a given column
func (r Row) Set(column string, value interface{}) {
	r[column] = value
}
