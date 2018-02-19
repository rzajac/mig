package mig

import (
    "strconv"
    "text/template"
    "time"
)

// Template function map.
var tplFnMap = template.FuncMap{
    "fTime": func(t time.Time) string { return strconv.FormatInt(t.UnixNano(), 10) },
}

// MySQL struct file template.
var mySQLStructTpl = template.Must(template.New("mig-mysql-struct-tpl").Funcs(tplFnMap).Parse(`package {{.Pkg}}

import "database/sql"

// A MigMySQL represents MySQL database migrations.
// The struct is used to attach methods with migrations.
type MigMySQL struct {
    db *sql.DB // Database handle.
}
`))

// MySQL migration file template.
var mySQLMigrationTpl = template.Must(template.New("mig-mysql-migration-tpl").Funcs(tplFnMap).Parse(`package {{.Pkg}}

import (
    "errors"
    "fmt"
)

// MigUp{{.Ts | fTime}}
func (m *MigMySQL) MigUp{{.Ts | fTime}}() error {
    fmt.Printf("up from {{.Ts | fTime}}\n")
    return errors.New("up migration {{.Ts | fTime}} not implemented")
}

// MigDown{{.Ts | fTime}}
func (m *MigMySQL) MigDown{{.Ts | fTime}}() error {
    fmt.Printf("down from {{.Ts | fTime}}\n")
    return errors.New("down migration {{.Ts | fTime}} not implemented")
}
`))
