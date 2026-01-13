package repl

import (
	"bufio"
	"fmt"
	"godb/engine"
	"godb/parser"
	"io"
	"strings"
)

// REPL represents the Read-Eval-Print Loop
type REPL struct {
	db     *engine.Database
	reader *bufio.Reader
}

// NewREPL creates a new REPL instance
func NewREPL(reader io.Reader) *REPL {
	return &REPL{
		db:     engine.NewDatabase(),
		reader: bufio.NewReader(reader),
	}
}

// Start begins the REPL loop
func (r *REPL) Start() {
	fmt.Println("godb - A minimal in-memory relational database")
	fmt.Println("Type 'exit' or 'quit' to exit")
	fmt.Println()

	for {
		fmt.Print("godb> ")

		// Read input
		input, err := r.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nGoodbye!")
				return
			}
			PrintError(fmt.Errorf("read error: %v", err))
			continue
		}

		input = strings.TrimSpace(input)

		// Handle empty input
		if input == "" {
			continue
		}

		// Handle exit commands
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			return
		}

		// Execute command
		r.executeCommand(input)
	}
}

// executeCommand parses and executes a command
func (r *REPL) executeCommand(input string) {
	// Parse command
	p := parser.NewParser(input)
	cmd, err := p.Parse()
	if err != nil {
		PrintError(fmt.Errorf("parse error: %v", err))
		return
	}

	// Execute based on command type
	switch c := cmd.(type) {
	case *parser.CreateTableCommand:
		r.executeCreateTable(c)
	case *parser.InsertCommand:
		r.executeInsert(c)
	case *parser.SelectCommand:
		r.executeSelect(c)
	case *parser.UpdateCommand:
		r.executeUpdate(c)
	case *parser.DeleteCommand:
		r.executeDelete(c)
	case *parser.JoinCommand:
		r.executeJoin(c)
	default:
		PrintError(fmt.Errorf("unknown command type"))
	}
}

// executeCreateTable executes a CREATE TABLE command
func (r *REPL) executeCreateTable(cmd *parser.CreateTableCommand) {
	err := r.db.CreateTable(cmd.TableName, cmd.Columns)
	if err != nil {
		PrintError(err)
		return
	}
	PrintSuccess(fmt.Sprintf("Table '%s' created successfully", cmd.TableName))
}

// executeInsert executes an INSERT command
func (r *REPL) executeInsert(cmd *parser.InsertCommand) {
	err := r.db.Insert(cmd.TableName, cmd.Values)
	if err != nil {
		PrintError(err)
		return
	}
	PrintSuccess("1 row inserted")
}

// executeSelect executes a SELECT command
func (r *REPL) executeSelect(cmd *parser.SelectCommand) {
	rows, err := r.db.Select(cmd.TableName, cmd.Columns, cmd.Condition)
	if err != nil {
		PrintError(err)
		return
	}
	PrintRows(rows)
}

// executeUpdate executes an UPDATE command
func (r *REPL) executeUpdate(cmd *parser.UpdateCommand) {
	count, err := r.db.Update(cmd.TableName, cmd.Updates, cmd.Condition)
	if err != nil {
		PrintError(err)
		return
	}
	PrintSuccess(fmt.Sprintf("%d row(s) updated", count))
}

// executeDelete executes a DELETE command
func (r *REPL) executeDelete(cmd *parser.DeleteCommand) {
	count, err := r.db.Delete(cmd.TableName, cmd.Condition)
	if err != nil {
		PrintError(err)
		return
	}
	PrintSuccess(fmt.Sprintf("%d row(s) deleted", count))
}

// executeJoin executes a JOIN command
func (r *REPL) executeJoin(cmd *parser.JoinCommand) {
	joinCondition := engine.JoinCondition{
		LeftColumn:  cmd.LeftColumn,
		RightColumn: cmd.RightColumn,
	}

	rows, err := r.db.InnerJoin(cmd.LeftTable, cmd.RightTable, joinCondition, cmd.SelectColumns)
	if err != nil {
		PrintError(err)
		return
	}
	PrintRows(rows)
}
