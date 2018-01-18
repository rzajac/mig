package mig

import (
    "fmt"
)

// Supported dialects
var dialects = [...]string{"mysql"}

// A Mig is a migrations manager.
type Mig struct {
    State   *Dir
    Dialect string
}

// NewMig creates new Mig instance.
func NewMig(root, dialect string) (*Mig, error) {
    if IsSupDialect(dialect) == false {
        return nil, fmt.Errorf("unsupported dialect: %s", dialect)
    }
    m := &Mig{Dialect: dialect}
    m.State = NewState(root)
    if m.State.Err != nil {
        return nil, m.State.Err
    }
    return m, nil
}

// Initialize initializes migration directory.
func (m *Mig) Initialize(dialect string) error {
    return m.State.Initialize(dialect)
}

// New creates new migration file for given dialect.
func (m *Mig) New(dialect string) (string, error) {
    return m.State.NewMigration(dialect)
}
