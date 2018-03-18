package mig

import (
    "testing"

    "github.com/spf13/afero"
    "github.com/stretchr/testify/assert"
)

func TestNewMig(t *testing.T) {
    // --- Given ---
    mem := afero.NewMemMapFs()
    fs = &afero.Afero{Fs: mem}
    fs.MkdirAll("a/b", 0755)
    afero.WriteFile(fs, "a/b/mig.yaml", []byte(tstYAML), 0644)

    // --- When ---
    cfg, err := NewMig("a/b/mig.yaml")

    // --- Then ---
    assert.NoError(t, err)
    assert.Exactly(t, "a/b/migrations", cfg.MigDir())
}

var tstYAML = `
dir: migrations
targets:
  setup:
    dialect: mysql
    dsn: setup_dsn
  auth:
    dialect: mysql
    dsn: auth_dsn
`
