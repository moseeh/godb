package parser

import (
	"fmt"
	"godb/engine"
	"strconv"
	"strings"
)

// Parser parses SQL commands from tokens
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser creates a new parser from input string
func NewParser(input string) *Parser {
	tokens := Tokenize(input)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse parses the input and returns a Command
func (p *Parser) Parse() (Command, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("empty input")
	}

	token := p.current()
	if token.Type != TokenKeyword {
		return nil, fmt.Errorf("expected keyword, got %s", token.Value)
	}

	keyword := strings.ToUpper(token.Value)
	switch keyword {
	case "CREATE":
		return p.parseCreateTable()
	case "INSERT":
		return p.parseInsert()
	case "SELECT":
		return p.parseSelect()
	case "UPDATE":
		return p.parseUpdate()
	case "DELETE":
		return p.parseDelete()
	default:
		return nil, fmt.Errorf("unknown command: %s", keyword)
	}
}

// parseCreateTable parses CREATE TABLE command
func (p *Parser) parseCreateTable() (*CreateTableCommand, error) {
	// CREATE TABLE table_name (col1 type [PRIMARY KEY], col2 type [UNIQUE], ...)
	p.advance() // Skip CREATE

	if !p.matchKeyword("TABLE") {
		return nil, fmt.Errorf("expected TABLE keyword")
	}
	p.advance()

	tableName, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	if !p.match(TokenLeftParen) {
		return nil, fmt.Errorf("expected '(' after table name")
	}
	p.advance()

	columns, err := p.parseColumnDefinitions()
	if err != nil {
		return nil, err
	}

	if !p.match(TokenRightParen) {
		return nil, fmt.Errorf("expected ')' after column definitions")
	}

	return &CreateTableCommand{
		TableName: tableName,
		Columns:   columns,
	}, nil
}

// parseColumnDefinitions parses column definitions in CREATE TABLE
func (p *Parser) parseColumnDefinitions() ([]engine.Column, error) {
	var columns []engine.Column

	for {
		colName, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}

		colTypeStr, err := p.expectKeyword()
		if err != nil {
			return nil, err
		}

		colType := engine.ColumnType(strings.ToUpper(colTypeStr))

		col := engine.Column{
			Name: colName,
			Type: colType,
		}

		// Check for PRIMARY KEY or UNIQUE
		for p.matchKeyword("PRIMARY") || p.matchKeyword("UNIQUE") || p.matchKeyword("NOT") {
			if p.matchKeyword("PRIMARY") {
				p.advance()
				if p.matchKeyword("KEY") {
					p.advance()
					col.PrimaryKey = true
					col.NotNull = true
				}
			} else if p.matchKeyword("UNIQUE") {
				p.advance()
				col.Unique = true
			} else if p.matchKeyword("NOT") {
				p.advance()
				if p.matchKeyword("NULL") {
					p.advance()
					col.NotNull = true
				}
			}
		}

		columns = append(columns, col)

		if p.match(TokenComma) {
			p.advance()
			continue
		}
		break
	}

	return columns, nil
}

// parseInsert parses INSERT INTO command
func (p *Parser) parseInsert() (*InsertCommand, error) {
	// INSERT INTO table_name VALUES (val1, val2, ...)
	p.advance() // Skip INSERT

	if !p.matchKeyword("INTO") {
		return nil, fmt.Errorf("expected INTO keyword")
	}
	p.advance()

	tableName, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	// Optional column list
	var columns []string
	if p.match(TokenLeftParen) {
		p.advance()
		columns, err = p.parseIdentifierList()
		if err != nil {
			return nil, err
		}
		if !p.match(TokenRightParen) {
			return nil, fmt.Errorf("expected ')' after column list")
		}
		p.advance()
	}

	if !p.matchKeyword("VALUES") {
		return nil, fmt.Errorf("expected VALUES keyword")
	}
	p.advance()

	if !p.match(TokenLeftParen) {
		return nil, fmt.Errorf("expected '(' after VALUES")
	}
	p.advance()

	values, err := p.parseValueList()
	if err != nil {
		return nil, err
	}

	if !p.match(TokenRightParen) {
		return nil, fmt.Errorf("expected ')' after values")
	}

	// Map values to columns
	row := make(engine.Row)
	if len(columns) > 0 {
		if len(columns) != len(values) {
			return nil, fmt.Errorf("column count doesn't match value count")
		}
		for i, col := range columns {
			row[col] = values[i]
		}
	} else {
		// If no columns specified, we can't proceed without schema
		return nil, fmt.Errorf("column names must be specified in INSERT")
	}

	return &InsertCommand{
		TableName: tableName,
		Values:    row,
	}, nil
}

// parseSelect parses SELECT command
func (p *Parser) parseSelect() (Command, error) {
	// SELECT col1, col2 FROM table [WHERE condition]
	// SELECT * FROM table1 INNER JOIN table2 ON table1.col = table2.col
	p.advance() // Skip SELECT

	columns, err := p.parseSelectColumns()
	if err != nil {
		return nil, err
	}

	if !p.matchKeyword("FROM") {
		return nil, fmt.Errorf("expected FROM keyword")
	}
	p.advance()

	tableName, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for JOIN
	if p.matchKeyword("INNER") {
		p.advance()
		if !p.matchKeyword("JOIN") {
			return nil, fmt.Errorf("expected JOIN after INNER")
		}
		p.advance()

		rightTable, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}

		if !p.matchKeyword("ON") {
			return nil, fmt.Errorf("expected ON after JOIN")
		}
		p.advance()

		// Parse join condition: table1.col = table2.col
		leftCol, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}

		if !p.matchOperator("=") {
			return nil, fmt.Errorf("expected '=' in JOIN condition")
		}
		p.advance()

		rightCol, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}

		// Extract column names without table prefix
		leftColName := extractColumnName(leftCol)
		rightColName := extractColumnName(rightCol)

		return &JoinCommand{
			LeftTable:     tableName,
			RightTable:    rightTable,
			LeftColumn:    leftColName,
			RightColumn:   rightColName,
			SelectColumns: columns,
		}, nil
	}

	// Parse WHERE clause if present
	var condition *engine.Condition
	if p.matchKeyword("WHERE") {
		p.advance()
		condition, err = p.parseCondition()
		if err != nil {
			return nil, err
		}
	}

	return &SelectCommand{
		TableName: tableName,
		Columns:   columns,
		Condition: condition,
	}, nil
}

// parseUpdate parses UPDATE command
func (p *Parser) parseUpdate() (*UpdateCommand, error) {
	// UPDATE table SET col1=val1, col2=val2 WHERE condition
	p.advance() // Skip UPDATE

	tableName, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	if !p.matchKeyword("SET") {
		return nil, fmt.Errorf("expected SET keyword")
	}
	p.advance()

	updates, err := p.parseSetClause()
	if err != nil {
		return nil, err
	}

	var condition *engine.Condition
	if p.matchKeyword("WHERE") {
		p.advance()
		condition, err = p.parseCondition()
		if err != nil {
			return nil, err
		}
	}

	return &UpdateCommand{
		TableName: tableName,
		Updates:   updates,
		Condition: condition,
	}, nil
}

// parseDelete parses DELETE command
func (p *Parser) parseDelete() (*DeleteCommand, error) {
	// DELETE FROM table WHERE condition
	p.advance() // Skip DELETE

	if !p.matchKeyword("FROM") {
		return nil, fmt.Errorf("expected FROM keyword")
	}
	p.advance()

	tableName, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	var condition *engine.Condition
	if p.matchKeyword("WHERE") {
		p.advance()
		condition, err = p.parseCondition()
		if err != nil {
			return nil, err
		}
	}

	return &DeleteCommand{
		TableName: tableName,
		Condition: condition,
	}, nil
}

// parseSelectColumns parses the column list in SELECT
func (p *Parser) parseSelectColumns() ([]string, error) {
	if p.current().Value == "*" {
		p.advance()
		return nil, nil // nil means all columns
	}

	return p.parseIdentifierList()
}

// parseIdentifierList parses a comma-separated list of identifiers
func (p *Parser) parseIdentifierList() ([]string, error) {
	var identifiers []string

	for {
		id, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, id)

		if p.match(TokenComma) {
			p.advance()
			continue
		}
		break
	}

	return identifiers, nil
}

// parseValueList parses a comma-separated list of values
func (p *Parser) parseValueList() ([]interface{}, error) {
	var values []interface{}

	for {
		val, err := p.expectValue()
		if err != nil {
			return nil, err
		}
		values = append(values, val)

		if p.match(TokenComma) {
			p.advance()
			continue
		}
		break
	}

	return values, nil
}

// parseSetClause parses SET col=val, col2=val2
func (p *Parser) parseSetClause() (engine.Row, error) {
	updates := make(engine.Row)

	for {
		col, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}

		if !p.matchOperator("=") {
			return nil, fmt.Errorf("expected '=' in SET clause")
		}
		p.advance()

		val, err := p.expectValue()
		if err != nil {
			return nil, err
		}

		updates[col] = val

		if p.match(TokenComma) {
			p.advance()
			continue
		}
		break
	}

	return updates, nil
}

// parseCondition parses a WHERE condition
func (p *Parser) parseCondition() (*engine.Condition, error) {
	col, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}

	if !p.match(TokenOperator) {
		return nil, fmt.Errorf("expected operator in condition")
	}
	op := p.current().Value
	p.advance()

	val, err := p.expectValue()
	if err != nil {
		return nil, err
	}

	return &engine.Condition{
		Column:   col,
		Operator: op,
		Value:    val,
	}, nil
}

// Helper functions

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) match(tokenType TokenType) bool {
	return p.current().Type == tokenType
}

func (p *Parser) matchKeyword(keyword string) bool {
	return p.match(TokenKeyword) && strings.EqualFold(p.current().Value, keyword)
}

func (p *Parser) matchOperator(op string) bool {
	return p.match(TokenOperator) && p.current().Value == op
}

func (p *Parser) expectIdentifier() (string, error) {
	if !p.match(TokenIdentifier) {
		return "", fmt.Errorf("expected identifier, got %v", p.current())
	}
	value := p.current().Value
	p.advance()
	return value, nil
}

func (p *Parser) expectKeyword() (string, error) {
	if !p.match(TokenKeyword) {
		return "", fmt.Errorf("expected keyword, got %v", p.current())
	}
	value := p.current().Value
	p.advance()
	return value, nil
}

func (p *Parser) expectValue() (interface{}, error) {
	token := p.current()

	switch token.Type {
	case TokenString:
		p.advance()
		return token.Value, nil
	case TokenNumber:
		p.advance()
		val, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", token.Value)
		}
		return val, nil
	case TokenKeyword:
		// Handle NULL, TRUE, FALSE
		upper := strings.ToUpper(token.Value)
		if upper == "NULL" {
			p.advance()
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected keyword in value position: %s", token.Value)
	default:
		return nil, fmt.Errorf("expected value, got %v", token)
	}
}

// extractColumnName extracts column name from qualified name (table.column)
func extractColumnName(qualified string) string {
	parts := strings.Split(qualified, ".")
	if len(parts) == 2 {
		return parts[1]
	}
	return qualified
}
