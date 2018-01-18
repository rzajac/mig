package mig

import (
    "fmt"
)

// Supported dialects
var dialects = [...]string{"mysql"}

// A Mig is a migrations manager.
type Mig struct {
    dir     *Dir
    dialect string
}

// NewMig creates new Mig instance.
func NewMig(root, dialect string) (*Mig, error) {
    if IsSupDialect(dialect) == false {
        return nil, fmt.Errorf("unsupported dialect: %s", dialect)
    }
    m := &Mig{dialect: dialect}
    m.dir = NewDir(root)
    if m.dir.Error != nil {
        return nil, m.dir.Error
    }
    return m, nil
}

// Init initializes migration directory.
func (m *Mig) Initialize() error {
    return m.dir.Init(m.dialect)
}

// New creates new migration file for given dialect.
func (m *Mig) New() (string, error) {
    return m.dir.NewMigration(m.dialect)
}
