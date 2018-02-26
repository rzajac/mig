package mig

import "sync"

// registry keeps all registered migrations.
var registry = newRegistry()

// reg represents migrations registered by Register.
type reg struct {
    sync.Mutex
    // Map with registered migrations.
    // The key is the name of the migration configuration.
    reg map[string][]Migration
}

// newRegistry creates new registry.
func newRegistry() *reg {
    return &reg{
        reg: make(map[string][]Migration),
    }
}

// Register registers migration.
func Register(name string, m Migration) {
    registry.Lock()
    defer registry.Unlock()
    registry.reg[name] = append(registry.reg[name], m)
}
