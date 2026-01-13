package engine

import "fmt"

// JoinCondition represents the condition for joining two tables
type JoinCondition struct {
	LeftColumn  string
	RightColumn string
}

// InnerJoin performs an INNER JOIN between two tables
func (db *Database) InnerJoin(leftTable, rightTable string, condition JoinCondition, selectColumns []string) ([]Row, error) {
	// Get both tables
	left, err := db.GetTable(leftTable)
	if err != nil {
		return nil, err
	}

	right, err := db.GetTable(rightTable)
	if err != nil {
		return nil, err
	}

	// Verify join columns exist
	if !left.hasColumn(condition.LeftColumn) {
		return nil, ErrColumnNotFound{
			TableName:  leftTable,
			ColumnName: condition.LeftColumn,
		}
	}

	if !right.hasColumn(condition.RightColumn) {
		return nil, ErrColumnNotFound{
			TableName:  rightTable,
			ColumnName: condition.RightColumn,
		}
	}

	var results []Row

	// Check if right table has an index on the join column
	rightIndex, hasIndex := right.GetIndex(condition.RightColumn)

	// Iterate through left table
	for _, leftRow := range left.rows {
		leftValue, ok := leftRow.Get(condition.LeftColumn)
		if !ok || leftValue == nil {
			continue
		}

		// Find matching rows in right table
		var matchingRightIndices []int
		if hasIndex {
			// Use index for faster lookup
			matchingRightIndices = rightIndex.Lookup(leftValue)
		} else {
			// Linear scan through right table
			for i, rightRow := range right.rows {
				rightValue, ok := rightRow.Get(condition.RightColumn)
				if ok && rightValue == leftValue {
					matchingRightIndices = append(matchingRightIndices, i)
				}
			}
		}

		// Create joined rows
		for _, rightIdx := range matchingRightIndices {
			if rightIdx >= len(right.rows) {
				continue
			}
			rightRow := right.rows[rightIdx]
			joinedRow := mergeRows(leftRow, rightRow, leftTable, rightTable)

			// Project columns if specified
			if len(selectColumns) > 0 {
				projectedRow := make(Row)
				for _, col := range selectColumns {
					if value, ok := joinedRow.Get(col); ok {
						projectedRow.Set(col, value)
					}
				}
				results = append(results, projectedRow)
			} else {
				results = append(results, joinedRow)
			}
		}
	}

	return results, nil
}

// mergeRows combines two rows from different tables, prefixing column names with table names
func mergeRows(left, right Row, leftTable, rightTable string) Row {
	result := make(Row)

	// Add columns from left table with prefix
	for col, val := range left {
		qualifiedCol := fmt.Sprintf("%s.%s", leftTable, col)
		result.Set(qualifiedCol, val)
	}

	// Add columns from right table with prefix
	for col, val := range right {
		qualifiedCol := fmt.Sprintf("%s.%s", rightTable, col)
		result.Set(qualifiedCol, val)
	}

	return result
}
