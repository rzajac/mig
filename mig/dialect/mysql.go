package dialect

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// MySQL represents mysql migration driver and connection.
type MySQL struct {
    *sql.DB
}

// NewMySQL returns new MySQL connection.
func NewMySQL(user, pass, host, schema string) (*MySQL, error) {
    var err error
    m := &MySQL{}
    dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8", user, pass, host, schema)
    if m.DB, err = sql.Open("mysql", dsn); err != nil {
        return nil, err
    }
    return m, nil
}

// MigEnsure checks if migration table is present and creates it if needed.
func (m *MySQL) MigEnsure() error {
    if _, err := m.Exec(mySQLMigTableCreate); err != nil {
        return err
    }
    return nil
}
