package mig

import "sort"

// Registered migrations.
var migrations = make(map[string][]Migration)

// Register registers Migration.
func Register(target string, migration Migration) {
    migrations[target] = append(migrations[target], migration)
}

// TargetNames returns registered target names.
func TargetNames() []string {
    names := make([]string, 0)
    for name := range migrations {
        names = append(names, name)
    }
    return names
}

// GetMigrations returns sorted migrations for given target.
func GetMigrations(target string) []Migration {
    migs, ok := migrations[target]
    if !ok {
        return make([]Migration, 0)
    }
    sort.Sort(migSort(migs))
    return migs
}

// Slice of Migrators with Sorter interface.
type migSort []Migration

func (m migSort) Len() int           { return len(m) }
func (m migSort) Less(i, j int) bool { return m[i].Version() < m[j].Version() }
func (m migSort) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
