package engine

import "fmt"

// ErrTableNotFound is returned when a table does not exist
type ErrTableNotFound struct {
	TableName string
}

func (e ErrTableNotFound) Error() string {
	return fmt.Sprintf("table '%s' does not exist", e.TableName)
}

// ErrTableAlreadyExists is returned when attempting to create a table that already exists
type ErrTableAlreadyExists struct {
	TableName string
}

func (e ErrTableAlreadyExists) Error() string {
	return fmt.Sprintf("table '%s' already exists", e.TableName)
}

// ErrPrimaryKeyViolation is returned when a primary key constraint is violated
type ErrPrimaryKeyViolation struct {
	TableName string
	Key       string
	Value     interface{}
}

func (e ErrPrimaryKeyViolation) Error() string {
	return fmt.Sprintf("primary key violation in table '%s': duplicate value '%v' for key '%s'",
		e.TableName, e.Value, e.Key)
}

// ErrUniqueViolation is returned when a unique constraint is violated
type ErrUniqueViolation struct {
	TableName string
	Column    string
	Value     interface{}
}

func (e ErrUniqueViolation) Error() string {
	return fmt.Sprintf("unique constraint violation in table '%s': duplicate value '%v' for column '%s'",
		e.TableName, e.Value, e.Column)
}

// ErrColumnNotFound is returned when a column does not exist in a table
type ErrColumnNotFound struct {
	TableName  string
	ColumnName string
}

func (e ErrColumnNotFound) Error() string {
	return fmt.Sprintf("column '%s' does not exist in table '%s'", e.ColumnName, e.TableName)
}

// ErrMissingRequiredColumn is returned when a required column is missing during insert/update
type ErrMissingRequiredColumn struct {
	TableName  string
	ColumnName string
}

func (e ErrMissingRequiredColumn) Error() string {
	return fmt.Sprintf("missing required column '%s' in table '%s'", e.ColumnName, e.TableName)
}

// ErrInvalidValue is returned when a value does not match expected type
type ErrInvalidValue struct {
	Column   string
	Expected string
	Got      interface{}
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value for column '%s': expected %s, got %T", e.Column, e.Expected, e.Got)
}

// ErrNoRowsAffected is returned when an update/delete operation affects no rows
type ErrNoRowsAffected struct{}

func (e ErrNoRowsAffected) Error() string {
	return "no rows affected by operation"
}

// ErrMultiplePrimaryKeys is returned when attempting to create a table with multiple primary keys
type ErrMultiplePrimaryKeys struct {
	TableName string
}

func (e ErrMultiplePrimaryKeys) Error() string {
	return fmt.Sprintf("table '%s' cannot have multiple primary keys", e.TableName)
}
