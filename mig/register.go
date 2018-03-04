package mig

import (
    "fmt"
    "sync"
)

// reg is a global variable which keeps all registered migrations.
var reg = &registry{
    reg: make(map[string][]Executor, 0),
}

// registry represents migrations registered by Register function.
type registry struct {
    sync.Mutex
    // Map with registered migrations.
    // The key is the name of the migration configuration.
    reg map[string][]Executor
}

// Register registers migration.
func Register(name string, m Executor) {
    reg.Lock()
    defer reg.Unlock()
    reg.reg[name] = append(reg.reg[name], m)
}

// List temporary debugging function.
// TODO
func List() {
    for n, m := range reg.reg {
        fmt.Println(n, m)
    }
}
