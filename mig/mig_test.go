package mig

import (
    "testing"

    "github.com/spf13/afero"
    "github.com/stretchr/testify/assert"
)

// prepareFs creates memory filesystem with mig.yaml file.
func prepareFs(configFile string) {
    // Set global filesystem abstraction to memory backed.
    fs = &afero.Afero{Fs: afero.NewMemMapFs()}
    fs.MkdirAll("a/b", 0755)
    afero.WriteFile(fs, "a/b/mig.yaml", []byte(configFile), 0644)
}

func TestMig_NewMig(t *testing.T) {
    // --- Given ---
    prepareFs(tstYAML)

    // --- When ---
    cfg, err := NewMig("a/b/mig.yaml")

    // --- Then ---
    assert.NoError(t, err)
    assert.Exactly(t, "a/b/migrations", cfg.MigDir())
    assert.Exactly(t, 2, len(cfg.Targets))
}

func TestMig_Target_t1(t *testing.T) {
    // --- Given ---
    prepareFs(tstYAML)

    // --- When ---
    cfg, _ := NewMig("a/b/mig.yaml")
    t1, err := cfg.Target("t1")

    // --- Then ---
    assert.NoError(t, err)
    assert.Exactly(t, "t1", t1.Name())
    assert.Exactly(t, "a/b/migrations/t1", t1.TargetDir())
}

var tstYAML = `
dir: migrations
targets:
  t1:
    dialect: mysql
    dsn: t1_dsn
  t2:
    dialect: bad
    dsn: t1_dsn
`

var td *tstDriver

func newTstDriver(name, dsn string) Driver {
    return td
}

type tstDriver struct {
    openReturn       error
    closeReturn      error
    initializeReturn error
    applyReturn      error
    revertReturn     error
    mergeReturn      error
    versionIntReturn int64
    versionReturn    error
    getMigByteReturn []byte
    getMigReturn     error
}

func (td *tstDriver) Open() error {
    return td.openReturn
}

func (td *tstDriver) Close() error {
    return td.closeReturn
}

func (td *tstDriver) Initialize() error {
    return td.initializeReturn
}

func (td *tstDriver) Apply(mig Migration) error {
    return td.applyReturn
}

func (td *tstDriver) Revert(mig Migration) error {
    return td.revertReturn
}

func (td *tstDriver) Merge(migs []Migration) error {
    return td.mergeReturn
}

func (td *tstDriver) Version() (int64, error) {
    return td.versionIntReturn, td.versionReturn
}

func (td *tstDriver) GenMigration(version int64) ([]byte, error) {
    return td.getMigByteReturn, td.getMigReturn
}
