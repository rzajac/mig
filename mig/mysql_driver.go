package mig

import (
    "database/sql"

    "github.com/go-sql-driver/mysql"
    "github.com/pkg/errors"
    "time"
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

func (m *mysqlDriver) Apply(migration Migrator) error {
    return nil
}

func (m *mysqlDriver) Revert(migration Migrator) error {
    return nil
}

func (m *mysqlDriver) Applied() (MigRows, error) {
    migs := make(MigRows, 0)
    rows, err := m.db.Query(mySQLGetApplied)
    switch {
    case err == sql.ErrNoRows:
        return migs, nil
    case err != nil:
        if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1146 {
            return nil, ErrNotInitialized
        } else {
            return nil, errors.WithStack(err)
        }
    }
    defer rows.Close()

    for rows.Next() {
        var v int64
        var t time.Time
        err := rows.Scan(&v, &t)
        if err != nil {
            return nil, err
        }
        migs[v] = t
    }
    if err := rows.Err(); err != nil {
        return nil, errors.WithStack(err)
    }
    return migs, nil
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

// Select applied migrations in descending order.
var mySQLGetApplied = `SELECT version, applied FROM migrations ORDER BY id ASC`
