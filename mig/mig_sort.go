package mig

// Slice of Migrators with Sorter interface.
type migSort []Migration

func (m migSort) Len() int           { return len(m) }
func (m migSort) Less(i, j int) bool { return m[i].Version() < m[j].Version() }
func (m migSort) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
