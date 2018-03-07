package mig

import (
    "fmt"
    "sort"
    "sync"
)

// Registered migrations.
var registry = &migrations{migs: make(map[string]migs, 0)}

// Register registers Migrator.
func Register(target string, mgr Migrator) {
    registry.Lock()
    defer registry.Unlock()
    registry.migs[target] = append(registry.migs[target], mgr)
}

// Slice of Migrators with Sorter interface.
type migs []Migrator

func (e migs) Len() int           { return len(e) }
func (e migs) Less(i, j int) bool { return e[i].Version() < e[j].Version() }
func (e migs) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// migrations represents list migrations.
type migrations struct {
    sync.Mutex
    migs map[string]migs
}

// sort sots all registered migrations.
func (m *migrations) sort() {
    m.Lock()
    defer m.Unlock()
    for _, m := range m.migs {
        sort.Sort(migs(m))
    }
}

// list lists migrations for given target.
func (m *migrations) list(target string) {
    for _, m := range registry.migs[target] {
        fmt.Println(target, m.Version())
    }
}
