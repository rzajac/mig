package mig

var migTplMainMysql = `
package migration

import "database/sql"

// A MigMySQL represents MySQL database migrations.
// The struct is used to attach methods with migrations.
type MigMySQL struct {
    schema  string  // Database name.
    Dialect string  // Migration Dialect.
    db      *sql.DB // Database handle.
}

// Dialect returns name of the Dialect the migration uses.
// The Dialect decides what driver is being passed to migration methods.
func (m *MigMySQL) Dialect() string {
    return m.Dialect
}
`

var migTplMigTpl = `
package migration

import "fmt"

// Mig%d
func (m *MigMySQL) Mig%d() error {
    fmt.Printf("dlasjdlkajhsdlkajhdslkajd\n")
    return nil
}
`
