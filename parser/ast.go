package parser

import "godb/engine"

// CommandType represents the type of SQL command
type CommandType int

const (
	CmdCreateTable CommandType = iota
	CmdInsert
	CmdSelect
	CmdUpdate
	CmdDelete
	CmdUnknown
)

// Command represents a parsed SQL command
type Command interface {
	Type() CommandType
}

// CreateTableCommand represents a CREATE TABLE statement
type CreateTableCommand struct {
	TableName string
	Columns   []engine.Column
}

func (c *CreateTableCommand) Type() CommandType {
	return CmdCreateTable
}

// InsertCommand represents an INSERT INTO statement
type InsertCommand struct {
	TableName string
	Values    engine.Row
}

func (c *InsertCommand) Type() CommandType {
	return CmdInsert
}

// SelectCommand represents a SELECT statement
type SelectCommand struct {
	TableName string
	Columns   []string
	Condition *engine.Condition
}

func (c *SelectCommand) Type() CommandType {
	return CmdSelect
}

// UpdateCommand represents an UPDATE statement
type UpdateCommand struct {
	TableName string
	Updates   engine.Row
	Condition *engine.Condition
}

func (c *UpdateCommand) Type() CommandType {
	return CmdUpdate
}

// DeleteCommand represents a DELETE statement
type DeleteCommand struct {
	TableName string
	Condition *engine.Condition
}

func (c *DeleteCommand) Type() CommandType {
	return CmdDelete
}

// JoinCommand represents a SELECT with INNER JOIN
type JoinCommand struct {
	LeftTable     string
	RightTable    string
	LeftColumn    string
	RightColumn   string
	SelectColumns []string
}

func (c *JoinCommand) Type() CommandType {
	return CmdSelect
}
