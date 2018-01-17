package mig

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestIsMigration(t *testing.T) {
    tt := []struct {
        name string
        exp  bool
    }{
        {"mig_mysql_1516144287290322398.go", true},
        {"mig_mysql.go", false},
        {"file.go", false},
    }

    for _, tc := range tt {
        assert.Exactly(t, tc.exp, IsMigration(tc.name))
    }
}

func TestDescMigration(t *testing.T) {
    tt := []struct {
        name    string
        dialect string
        ts      string
        err     bool
    }{
        {"mig_mysql_1516144287290322398.go", "mysql", "1516144287290322398", false},
    }

    for _, tc := range tt {
        dialect, ts, err := FileDialectAndTs(tc.name)
        if tc.err {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
        }

        assert.Exactly(t, tc.dialect, dialect)
        assert.Exactly(t, tc.ts, ts)
    }
}
