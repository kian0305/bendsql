package perf

type OutputMetaData struct {
	Table string `json:"table"`
	Tag   string `json:"tag"`
	Size  string `json:"size"`
}

type OutputSchema struct {
	Name     string    `json:"name"`
	SQL      string    `json:"sql"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Median   float64   `json:"median"`
	StdDev   float64   `json:"std_dev"`
	ReadRow  uint64    `json:"read_row"`
	ReadByte uint64    `json:"read_byte"`
	Time     []float64 `json:"time"`
	Error    []string  `json:"error"`
	Mean     float64   `json:"mean"`
}

type OutputFile struct {
	MetaData OutputMetaData `json:"metadata"`
	Schema   []OutputSchema `json:"schema"`
}
