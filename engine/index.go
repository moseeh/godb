package engine

// Index represents a hash-based index for a column
// Maps column value -> list of row indices
type Index struct {
	column string
	data   map[interface{}][]int
}

// NewIndex creates a new index for a column
func NewIndex(column string) *Index {
	return &Index{
		column: column,
		data:   make(map[interface{}][]int),
	}
}

// Add adds a row index to the index for a given value
func (idx *Index) Add(value interface{}, rowIndex int) {
	if value == nil {
		return // Don't index nil values
	}
	idx.data[value] = append(idx.data[value], rowIndex)
}

// Remove removes a row index from the index for a given value
func (idx *Index) Remove(value interface{}, rowIndex int) {
	if value == nil {
		return
	}

	indices, exists := idx.data[value]
	if !exists {
		return
	}

	// Filter out the row index
	newIndices := make([]int, 0, len(indices))
	for _, i := range indices {
		if i != rowIndex {
			newIndices = append(newIndices, i)
		}
	}

	if len(newIndices) == 0 {
		delete(idx.data, value)
	} else {
		idx.data[value] = newIndices
	}
}

// Lookup returns all row indices that match the given value
func (idx *Index) Lookup(value interface{}) []int {
	if value == nil {
		return nil
	}
	return idx.data[value]
}

// Update updates the index when a row's value changes
func (idx *Index) Update(oldValue, newValue interface{}, rowIndex int) {
	idx.Remove(oldValue, rowIndex)
	idx.Add(newValue, rowIndex)
}

// Has checks if a value exists in the index
func (idx *Index) Has(value interface{}) bool {
	if value == nil {
		return false
	}
	_, exists := idx.data[value]
	return exists
}
