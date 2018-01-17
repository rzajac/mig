package mig

var migTplMainMysql = `
package migration

import "database/sql"

// A MigMySQL represents MySQL database migrations.
// The struct is used to attach methods with migrations.
type MigMySQL struct {
    schema  string  // Database name.
    dialect string  // Migration dialect.
    db      *sql.DB // Database handle.
}

// Dialect returns name of the dialect the migration uses.
// The dialect decides what driver is being passed to migration methods.
func (m *MigMySQL) Dialect() string {
    return m.dialect
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