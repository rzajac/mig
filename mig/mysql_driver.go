package mig

import (
    "database/sql"
    "time"

    "github.com/go-sql-driver/mysql"
    "github.com/pkg/errors"
)

// mysqlDriver represents MySQL migration driver.
type mysqlDriver struct {
    target Target
    db     *sql.DB
}

// newMYSQLDriver returns new instance of mysqlDriver.
func newMYSQLDriver(config Target) *mysqlDriver {
    return &mysqlDriver{target: config}
}

func (m *mysqlDriver) Open() error {
    if m.db != nil {
        return nil
    }
    var err error
    m.db, err = sql.Open("mysql", m.target.Dsn())
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

func (m *mysqlDriver) Creator() Creator {
    return &mysqlCreator{target: m.target}
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
