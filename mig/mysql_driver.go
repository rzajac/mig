package mig

import (
    "database/sql"

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

func (m *mysqlDriver) Apply(migration Migrator) error {
    return nil
}

func (m *mysqlDriver) Revert(migration Migrator) error {
    return nil
}

func (m *mysqlDriver) Applied() ([]int64, error) {
    var vers []int64
    rows, err := m.db.Query(mySQLGetApplied)
    switch {
    case err == sql.ErrNoRows:
        return vers, nil
    case err != nil:
        return nil, errors.WithStack(err)
    }
    defer rows.Close()

    for rows.Next() {
        var v int64
        err := rows.Scan(&version, &r.info, &r.createdAt)
        if err != nil {
            return nil, err
        }
        if len(desc) == 0 {
            r.current = true
        }
        desc = append(desc, &r)
    }
    if err := rows.Err(); err != nil {
        return nil, errors.WithStack(err)
    }
    return desc, nil
}

func (m *mysqlDriver) Initialize() error {
    if _, err := m.db.Exec(mySQLMigTableCreate); err != nil {
        return err
    }
    return nil
}

func (m *mysqlDriver) Creator() Creator {
    return newMySQLCreator(m.target)
}

// Create migrations table.
var mySQLMigTableCreate = `CREATE TABLE migrations (
  version BIGINT UNSIGNED NOT NULL,
  applied TIMESTAMP NOT NULL,
  PRIMARY KEY (version)
) ENGINE=InnoDB`

// Select most recent migration version.
var mySQLGetVersion = `SELECT version FROM migrations ORDER BY version DESC LIMIT 1`
// Select applied migrations in descending order.
var mySQLGetApplied = `SELECT version, applied FROM migrations ORDER BY id ASC`
