package mig

// Slice of Migrators with Sorter interface.
type migSort []Migration

func (e migSort) Len() int           { return len(e) }
func (e migSort) Less(i, j int) bool { return e[i].Version() < e[j].Version() }
func (e migSort) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
