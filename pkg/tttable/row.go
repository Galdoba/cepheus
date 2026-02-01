package tttable

// Row represents a single row in the table
type Row struct {
	Key       string   `json:"key" toml:"key"  yaml:"key"`
	Value     string   `json:"value" toml:"value"  yaml:"value"`
	NextTable string   `json:"next_table,omitempty" toml:"next_table,omitempty"  yaml:"next_table,omitempty"`
	ApplyMods []string `json:"apply_mods,omitempty" toml:"apply_mods,omitempty"  yaml:"apply_mods,omitempty"`
	RollAgain bool     `json:"roll_again,omitempty" toml:"roll_again,omitempty"  yaml:"roll_again,omitempty"`
	ReRoll    bool     `json:"re_roll,omitempty" toml:"re_roll,omitempty"  yaml:"re_roll,omitempty"`
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
