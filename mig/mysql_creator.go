package mig

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "path"
    "text/template"

    "github.com/pkg/errors"
)

// mysqlCreator represents MySQL migration file creator.
type mysqlCreator struct {
    target Target
}

// newMySQLCreator returns new mysqlCreator based on target.
func newMySQLCreator(target Target) *mysqlCreator {
    return &mysqlCreator{
        target: target,
    }
}

func (cr *mysqlCreator) CreateMigration(version int64) error {
    if err := cr.ensure(); err != nil {
        return err
    }
    var data = struct {
        Pkg     string
        Version int64
    }{
        Pkg:     cr.target.Name(),
        Version: version,
    }
    var buf bytes.Buffer
    if err := mySQLMigTpl.Execute(&buf, data); err != nil {
        return errors.WithStack(err)
    }
    if err := ioutil.WriteFile(cr.path(version), buf.Bytes(), 0666); err != nil {
        return errors.WithStack(err)
    }
    return nil
}

// ensureFs ensures everything is ready to create migration files.
func (cr *mysqlCreator) ensure() error {
    if err := checkCreateDir(cr.target.MigDir()); err != nil {
        return err
    }
    return nil
}

// path returns absolute path to migration file with given version.
func (cr *mysqlCreator) path(version int64) string {
    return path.Join(cr.target.MigDir(), fmt.Sprintf("%d.go", version))
}

// MySQL migration file template.
var mySQLMigTpl = template.Must(template.New("registry-mysqlDriver-struct-tpl").Parse(`package {{.Pkg}}

import (
    "database/sql"
    "time"

    "github.com/rzajac/mig/mig"
)

func init() {
    registry.Register("setup", &Mig{{.Version}}{})
}

type Mig{{.Version}} struct {
    createdAt time.Time
    db        *sql.DB
}

func (m *Mig{{.Version}}) Setup(driver interface{}, createdAt time.Time) {
    m.createdAt = createdAt
    m.db = driver.(*sql.DB)
}

func (m *Mig{{.Version}}) Version() int64 {
    return {{.Version}}
}

func (m *Mig{{.Version}}) CreatedAt() time.Time {
    return m.createdAt
}

// ======================= DO NOT EDIT ABOVE THIS LINE =======================

func (m *Mig{{.Version}}) Apply() error {
    _, err := m.db.Exec("")
    return err
}

func (m *Mig{{.Version}}) Revert() error {
    _, err := m.db.Exec("")
    return err
}

func (m *Mig{{.Version}}) Info() string {
    return "example description for version {{.Version}}"
}
`))
