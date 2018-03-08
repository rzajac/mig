package mig

import (
    "fmt"
    "sort"
    "sync"
    "time"

    "github.com/pkg/errors"
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

// applyFromDb apply creation time and driver to migrations for given target.
func (m *migrations) applyFromDb(drv Driver, target string, rows MigRows) error {
    if len(rows) > len(m.migs[target]) {
        return errors.New("database has more migrations then filesystem")
    }
    for _, mgr := range m.migs[target] {
        ver := mgr.Version()
        t, ok := rows[ver]
        if !ok {
            t = time.Time{}
        }
        mgr.Setup(drv.Drv(), t)
    }
    m.sort()
    if err := m.validate(target); err != nil {
        return err
    }
    return nil
}

// validate validates migrations list for given target.
func (m *migrations) validate(target string) error {
    // No migrations no possibility for error.
    if len(m.migs[target]) == 0 {
        return nil
    }
    prev := m.migs[target][0].AppliedAt().IsZero()
    for _, mgr := range m.migs[target] {
        curr := mgr.AppliedAt().IsZero()
        switch {
        case prev == false && curr == true:
            return errors.New("migrations are not continuous")
        default:
            prev = curr
        }
    }
    return nil
}

// sort sots all registered migrations.
func (m *migrations) sort() {
    m.Lock()
    defer m.Unlock()
    for _, m := range m.migs {
        sort.Sort(migs(m))
    }
}

func (m *migrations) applyAll(target string) error {
    for _, mgr := range m.migs[target] {
        if !mgr.AppliedAt().IsZero() {
            continue
        }
        if err := mgr.Apply(); err != nil {
            return err
        }
    }
    return nil
}

// list lists migrations for given target.
func (m *migrations) list(target string) {
    for _, m := range registry.migs[target] {
        fmt.Println(target, m.Version())
    }
}
