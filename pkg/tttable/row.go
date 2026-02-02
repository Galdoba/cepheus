package tttable

// Row represents a single row in the table
type TableEntry struct {
	Key            string   `json:"key,omitempty" toml:"key,omitempty" yaml:"key,omitempty"`
	Value          string   `json:"value,omitempty" toml:"value,omitempty" yaml:"value,omitempty"`
	NextTable      string   `json:"next_table,omitempty" toml:"next_table,omitempty" yaml:"next_table,omitempty"`
	ModsFlat       []string `json:"mods_flat,omitempty" toml:"mods_flat,omitempty" yaml:"mods_flat,omitempty"`
	ModsCumulative []string `json:"mods_cumulative,omitempty" toml:"mods_cumulative,omitempty" yaml:"mods_cumulative,omitempty"`
	ModsMax        []string `json:"mods_max,omitempty" toml:"mods_max,omitempty" yaml:"mods_max,omitempty"`
	ModsMin        []string `json:"mods_min,omitempty" toml:"mods_min,omitempty" yaml:"mods_min,omitempty"`
	RollAgain      bool     `json:"roll_again,omitempty" toml:"roll_again,omitempty" yaml:"roll_again,omitempty"`
	ReRoll         bool     `json:"re_roll,omitempty" toml:"re_roll,omitempty" yaml:"re_roll,omitempty"`
}

// NewRow creates a new Row with given key and value
func NewTableEntry(key, value string) TableEntry {
	return TableEntry{
		Key:   key,
		Value: value,
	}
}

// RowOption is a functional option for configuring a TableEntry
type TableEntryOption func(*TableEntry) error
