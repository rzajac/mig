package migration

import (
    "database/sql"
    "errors"
)

// A Migration represents database migrations.
// The struct is used to attach methods with migrations.
type Migration struct{
    // Database name.
    schema string
}

// Dialect returns name of the dialect the migration uses.
// The dialect decides what driver is being passed to migration methods.
func (m *Migration) Dialect() string {
    return "mysql"
}

// MigTableName returns name of the migration table.
func (m *Migration) TableName() string {
    return "migrations"
}

// Check checks if migration table is present and creates it if needed.
func (m *Migration) Check(db *sql.DB) error {
    var exist bool
    rows, err := db.Query("SELECT `table_name` FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `table_schema` = ?", m.schema)
    if err != nil {
        return err
    }
    var tn string
    for rows.Next() {
        if err := rows.Scan(&tn); err != nil {
            return err
        }
        if tn == m.TableName() {
            exist = true
            break
        }
    }
    rows.Close()

    if !exist {
        // Create migrations table.
    }

    return errors.New("check is not implemented")
}
