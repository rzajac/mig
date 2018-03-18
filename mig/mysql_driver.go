package mig

import (
    "bytes"
    "database/sql"
    "text/template"
    "time"

    "github.com/go-sql-driver/mysql"
    "github.com/pkg/errors"
)

// mysqlDriver represents MySQL migration driver.
type mysqlDriver struct {
    name string
    dsn  string
    db   *sql.DB
}

// newMYSQLDriver returns new instance of mysqlDriver.
func newMYSQLDriver(name, dsn string) *mysqlDriver {
    return &mysqlDriver{name: name, dsn: dsn}
}

func (m *mysqlDriver) Open() error {
    if m.db != nil {
        return nil
    }
    var err error
    m.db, err = sql.Open("mysql", m.dsn)
    return err
}

func (m *mysqlDriver) Close() error {
    return m.db.Close()
}

func (m *mysqlDriver) Version() (int64, error) {
    var v int64
    row := m.db.QueryRow(mySQLGetVersion)
    err := row.Scan(&v)
    if err == sql.ErrNoRows {
        return 0, nil
    }
    if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1146 {
        return 0, ErrNotInitialized
    }
    return v, err
}

func (m *mysqlDriver) Apply(migration Migration) error {
    return nil
}

func (m *mysqlDriver) Revert(migration Migration) error {
    return nil
}

func (m *mysqlDriver) Merge(migs []Migration) error {
    ml := make(map[int64]time.Time, 0)
    rows, err := m.db.Query(mySQLGetApplied)
    switch {
    case err == sql.ErrNoRows:
        return nil
    case err != nil:
        if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1146 {
            return ErrNotInitialized
        } else {
            return errors.WithStack(err)
        }
    }
    defer rows.Close()

    for rows.Next() {
        var t time.Time
        var v int64
        err := rows.Scan(&v, &t)
        if err != nil {
            return err
        }
        ml[v] = t
    }
    if err := rows.Err(); err != nil {
        return errors.WithStack(err)
    }
    for _, mig := range migs {
        t, ok := ml[mig.Version()]
        if !ok {
            t = time.Time{}
        }
        mig.Setup(m.db, t)
    }
    return nil
}

func (m *mysqlDriver) Initialize() error {
    if _, err := m.db.Exec(mySQLMigTableCreate); err != nil {
        return err
    }
    return nil
}

func (m *mysqlDriver) GenMigration(version int64) ([]byte, error) {
    var data = struct {
        Pkg     string
        Version int64
    }{
        Pkg:     m.name,
        Version: version,
    }

    var buf bytes.Buffer
    if err := mySQLMigTpl.Execute(&buf, data); err != nil {
        return []byte{}, errors.WithStack(err)
    }
    return buf.Bytes(), nil
}

// Create migrations table.
var mySQLMigTableCreate = `CREATE TABLE migrations (
  version BIGINT UNSIGNED NOT NULL,
  applied TIMESTAMP NOT NULL,
  PRIMARY KEY (version)
) ENGINE=InnoDB`

// Select applied migrations in descending order.
var mySQLGetApplied = `SELECT version, applied FROM migrations ORDER BY version ASC`
// Select most recent migration version.
var mySQLGetVersion = `SELECT version FROM migrations ORDER BY version DESC LIMIT 1`

// MySQL migration file template.
var mySQLMigTpl = template.Must(template.New("mysql-mig-tpl").Parse(`package {{.Pkg}}

import (
    "database/sql"
    "time"

    "github.com/rzajac/mig/mig"
)

func init() {
    mig.Register("{{.Pkg}}", &Mig{{.Version}}{})
}

type Mig{{.Version}} struct {
    appliedAt time.Time
    db        *sql.DB
}

func (m *Mig{{.Version}}) Setup(driver interface{}, appliedAt time.Time) {
    m.appliedAt = appliedAt
    m.db = driver.(*sql.DB)
}

func (m *Mig{{.Version}}) Version() int64 {
    return {{.Version}}
}

func (m *Mig{{.Version}}) AppliedAt() time.Time {
    return m.appliedAt
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
