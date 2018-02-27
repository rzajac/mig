package mig

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "path"
    "text/template"

    "github.com/pkg/errors"
)

// mysqlCreator represents MySQL migrations creator.
type mysqlCreator struct {
    cfg DbConfigurator
}

// NewMySQLCreator returns new MySQL migration creator.
func NewMySQLCreator(config DbConfigurator) *mysqlCreator {
    return &mysqlCreator{
        cfg: config,
    }
}

// CreateMigration creates new migration file.
func (c *mysqlCreator) CreateMigration(version int64) error {
    if err := c.ensure(); err != nil {
        return err
    }
    var data = struct {
        Pkg     string
        Version int64
    }{
        Pkg:     c.cfg.Name(),
        Version: version,
    }
    var buf bytes.Buffer
    if err := mySQLMigTpl.Execute(&buf, data); err != nil {
        return errors.WithStack(err)
    }
    if err := ioutil.WriteFile(c.migFilePath(version), buf.Bytes(), 0666); err != nil {
        return errors.WithStack(err)
    }
    return nil
}

// ensure everything is ready to create migration files.
func (c *mysqlCreator) ensure() error {
    if err := checkCreateDir(c.cfg.MigDir()); err != nil {
        return err
    }
    return nil
}

// migFilePath returns absolute path to migration file with given version.
func (c *mysqlCreator) migFilePath(version int64) string {
    return path.Join(c.cfg.MigDir(), fmt.Sprintf("%d.go", version))
}

// MySQL migration file template.
var mySQLMigTpl = template.Must(template.New("mig-mysqlDriver-struct-tpl").Parse(`package {{.Pkg}}

import (
    "database/sql"

    "github.com/rzajac/mig/mig"
)

func init() {
    mig.Register("setup", &Mig{{.Version}}{})
}

type Mig{{.Version}} struct {
    applied bool
    db      *sql.DB
}

func (m *Mig{{.Version}}) Setup(driver interface{}, applied bool) {
    m.applied = applied
    m.db = driver.(*sql.DB)
}

func (m *Mig{{.Version}}) Version() int64 {
    return {{.Version}}
}

func (m *Mig{{.Version}}) IsApplied() bool {
    return m.applied
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

func (m *Mig{{.Version}}) Description() string {
    return "example description for version {{.Version}}"
}
`))
