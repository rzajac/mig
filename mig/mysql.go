package mig

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// NewMySQL returns new MySQL connection.
func NewMySQL(user, pass, host, schema string) *sql.DB {
    dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8", user, pass, host, schema)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    return db
}
