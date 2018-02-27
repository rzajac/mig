package mig

import (
    "fmt"
    "sync"
)

// registry keeps all registered migrations.
var registry = &reg{
    reg: make(map[string][]Migration),
}

// reg represents migrations registered by Register.
type reg struct {
    sync.Mutex
    // Map with registered migrations.
    // The key is the name of the migration configuration.
    reg map[string][]Migration
}

// Register registers migration.
func Register(name string, m Migration) {
    registry.Lock()
    defer registry.Unlock()
    registry.reg[name] = append(registry.reg[name], m)
}

func List() {
    for n, m := range registry.reg {
        fmt.Println(n, m)
    }
}
