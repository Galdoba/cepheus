package tttable

// Row represents a single row in the table
type Row struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewRow creates a new Row with given key and value
func NewRow(key, value string) Row {
	return Row{
		Key:   key,
		Value: value,
	}
}

// RowOption is a functional option for configuring a Row
type RowOption func(*Row)

// WithKey sets the key for a row
func WithKey(key string) RowOption {
	return func(r *Row) {
		r.Key = key
	}
}

// WithValue sets the value for a row
func WithValue(value string) RowOption {
	return func(r *Row) {
		r.Value = value
	}
}

// NewRowWithOptions creates a new Row with functional options
func NewRowWithOptions(opts ...RowOption) Row {
	r := &Row{}
	for _, opt := range opts {
		opt(r)
	}
	return *r
}
