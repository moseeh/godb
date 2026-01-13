package engine

// ColumnType represents the data type of a column
type ColumnType string

const (
	TypeInt    ColumnType = "INT"
	TypeString ColumnType = "STRING"
	TypeBool   ColumnType = "BOOL"
)

// Column represents a table column with its schema
type Column struct {
	Name       string
	Type       ColumnType
	PrimaryKey bool
	Unique     bool
	NotNull    bool
}

// ConstraintChecker validates constraints on rows
type ConstraintChecker struct {
	table *Table
}

// NewConstraintChecker creates a new constraint checker for a table
func NewConstraintChecker(table *Table) *ConstraintChecker {
	return &ConstraintChecker{table: table}
}

// ValidateInsert checks if a row can be inserted without violating constraints
func (c *ConstraintChecker) ValidateInsert(row Row) error {
	// Check primary key constraint
	if c.table.primaryKey != "" {
		pkValue, hasPK := row.Get(c.table.primaryKey)
		if !hasPK {
			return ErrMissingRequiredColumn{
				TableName:  c.table.name,
				ColumnName: c.table.primaryKey,
			}
		}

		// Check for duplicate primary key
		if c.table.hasPrimaryKeyValue(pkValue) {
			return ErrPrimaryKeyViolation{
				TableName: c.table.name,
				Key:       c.table.primaryKey,
				Value:     pkValue,
			}
		}
	}

	// Check unique constraints
	for _, col := range c.table.schema {
		if col.Unique && col.Name != c.table.primaryKey {
			value, hasValue := row.Get(col.Name)
			if hasValue && c.table.hasUniqueValue(col.Name, value) {
				return ErrUniqueViolation{
					TableName: c.table.name,
					Column:    col.Name,
					Value:     value,
				}
			}
		}
	}

	// Check not null constraints
	for _, col := range c.table.schema {
		if col.NotNull || col.PrimaryKey {
			value, hasValue := row.Get(col.Name)
			if !hasValue || value == nil {
				return ErrMissingRequiredColumn{
					TableName:  c.table.name,
					ColumnName: col.Name,
				}
			}
		}
	}

	return nil
}

// ValidateUpdate checks if a row can be updated without violating constraints
func (c *ConstraintChecker) ValidateUpdate(oldRow, newRow Row) error {
	// Check primary key constraint (if primary key is being changed)
	if c.table.primaryKey != "" {
		oldPK, _ := oldRow.Get(c.table.primaryKey)
		newPK, _ := newRow.Get(c.table.primaryKey)

		// If primary key changed, check for duplicates
		if oldPK != newPK && c.table.hasPrimaryKeyValue(newPK) {
			return ErrPrimaryKeyViolation{
				TableName: c.table.name,
				Key:       c.table.primaryKey,
				Value:     newPK,
			}
		}
	}

	// Check unique constraints
	for _, col := range c.table.schema {
		if col.Unique && col.Name != c.table.primaryKey {
			oldValue, _ := oldRow.Get(col.Name)
			newValue, hasNewValue := newRow.Get(col.Name)

			// If value changed, check for duplicates
			if hasNewValue && oldValue != newValue && c.table.hasUniqueValue(col.Name, newValue) {
				return ErrUniqueViolation{
					TableName: c.table.name,
					Column:    col.Name,
					Value:     newValue,
				}
			}
		}
	}

	// Check not null constraints
	for _, col := range c.table.schema {
		if col.NotNull || col.PrimaryKey {
			value, hasValue := newRow.Get(col.Name)
			if !hasValue || value == nil {
				return ErrMissingRequiredColumn{
					TableName:  c.table.name,
					ColumnName: col.Name,
				}
			}
		}
	}

	return nil
}
