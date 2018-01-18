package dialect

import (
    tpl "text/template"
    "time"
    "strconv"
)

var tplFuncs = tpl.FuncMap{
    "fTime": func(t time.Time) string { return strconv.FormatInt(t.UnixNano(), 10) },
}

var MysqlDialectTpl = tpl.Must(tpl.New("mig-mysql-dialect-tpl").Funcs(tplFuncs).Parse(`
package migration

import "database/sql"

// A MigMySQL represents MySQL database migrations.
// The struct is used to attach methods with migrations.
type MigMySQL struct {
    db      *sql.DB // Database handle.
}
`))

var MysqlMigTpl = tpl.Must(tpl.New("mig-mysql-migration-tpl").Funcs(tplFuncs).Parse(`
package migration

import (
    "errors"
    "fmt"
)

// MigUp{{.Ts | fTime}}
func (m *MigMySQL) Mig{{.Ts | fTime}}_up() error {
    fmt.Printf("up from {{.Ts | fTime}}\n")
    return errors.New("up migration {{.Ts | fTime}} not implemented")
}

// MigDown{{.Ts | fTime}}
func (m *MigMySQL) Mig{{.Ts | fTime}}_down() error {
    fmt.Printf("down from {{.Ts | fTime}}\n")
    return errors.New("down migration {{.Ts | fTime}} not implemented")
}
`))
