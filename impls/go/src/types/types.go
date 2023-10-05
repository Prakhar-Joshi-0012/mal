package types

type MalType interface {
}

// MalList
type MalList struct {
	Val  []MalType
	Meta MalType
}

// MalSymbols
type Malsymbols struct {
	Val  string
	Meta MalType
}

// Atoms
type Atoms struct {
	Val  MalType
	Meta MalType
}

// Vectors
type MalVector struct {
	Val  []MalType
	Meta MalType
}
type MalHash struct {
	Val  []MalType
	Meta MalType
}
