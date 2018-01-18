package mig

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
)

// A Dir represents migrations directory state.
type Dir struct {
    // The absolute path to directory containing migration directory.
    RootDir string
    // RootDir is a directory
    RootIsDir bool
    // The absolute path to migration directory.
    MigDir string
    // MigDir is a directory.
    MigIsDir bool
    // Path has migration directory.
    HasMigDir bool
    // Set to true if migration directory is empty.
    MigDirEmpty bool
    // Set to true if migration directory exists
    // and is initialized with at least one Dialect.
    Initialized bool
    // Set to true when directory pointed by path is ready to use by mig.
    Valid bool
    // List of initialized dialects.
    Dialects []string
    // The list of migrations.
    Migrations []string
    // The error which stopped the check.
    Err error
}

// NewState creates new Dir instance and collects all the needed info.
func NewState(root string) *Dir {
    s := &Dir{MigDirEmpty: true}
    s.RootDir, s.Err = filepath.Abs(root)
    if s.Err != nil {
        return s
    }
    s.Collect()
    return s
}

// Initialize initialize migrations directory for given Dialect.
func (d *Dir) Initialize(dialect string) error {
    if IsSupDialect(dialect) == false {
        return d.setErr(fmt.Errorf("unsupported dialect: %d", dialect))
    }
    if d.InitializedFor(dialect) == true {
        return nil
    }
    if d.MigIsDir == false {
        if err := d.CreateMigDir(); err != nil {
            return d.setErr(err)
        }
    }
    fn := path.Join(d.RootDir, MigStructFileName(dialect))
    if err := ioutil.WriteFile(fn, []byte(migTplMainMysql), 0666); err != nil {
        return d.setErr(err)
    }
    return nil
}

// NewMigration creates new migration file for given Dialect.
func (d *Dir) NewMigration(dialect string) (string, error) {
    if IsSupDialect(dialect) == false {
        return "", d.setErr(fmt.Errorf("unsupported dialect: %d", dialect))
    }
    ts, fn := MigFileName(dialect)
    fn = path.Join(d.RootDir, fn)
    data := []byte(fmt.Sprintf(migTplMigTpl, ts, ts))

    if err := ioutil.WriteFile(fn, data, 0666); err != nil {
        return "", d.setErr(err)
    }
    return fn, d.Collect()
}

func (d *Dir) CreateMigDir() error {
    if err := CreateDir(d.MigDir); err != nil {
        d.Err = err
        d.Valid = false
        return err
    }
    d.Collect()
    return nil
}

// Collect collects all the information and sets all fields accordingly.
func (d *Dir) Collect() error {
    var err error
    if d.RootIsDir, err = IsDir(d.RootDir); err != nil {
        return d.setErr(err)
    }

    if d.RootIsDir == false {
        return d.setErr(fmt.Errorf("%d must be an empty directory", d.RootDir))
    }

    // Check migration directory exists in root and it'd actually a directory.
    var ex bool
    d.MigDir = filepath.Join(d.RootDir, "migration")
    if ex, err = FileExists(d.MigDir); err != nil {
        return d.setErr(err)
    }
    if d.MigIsDir, err = IsDir(d.MigDir); err != nil {
        return d.setErr(err)
    }
    if ex && !d.MigIsDir {
        return d.setErr(fmt.Errorf("%d is not a directory", d.MigDir))
    }

    if d.HasMigDir == false {
        d.Valid = true
        return nil
    }

    var fs []os.FileInfo
    if fs, err = ioutil.ReadDir(d.MigDir); err != nil {
        return d.setErr(err)
    }
    for _, fi := range fs {
        switch {
        case IsMigStructFile(fi.Name()):
            dialect, _ := StructFileParts(fi.Name())
            if !IsSupDialect(dialect) {
                return d.setErr(fmt.Errorf("unsupported Dialect %d in migration file %d", dialect, fi.Name()))
            }
            d.Dialects = append(d.Dialects, dialect)
        case IsMigFile(fi.Name()):
            d.Migrations = append(d.Migrations, fi.Name())
        default:
            return d.setErr(fmt.Errorf("unexpected file %d", fi.Name()))
        }
    }

    if len(d.Dialects) > 0 || len(d.Migrations) > 0 {
        d.MigDirEmpty = false
        d.Initialized = true
    }
    d.Valid = true
    return nil
}

// InitializedFor returns true if migrations have been
// initialized for given dialect.
func (d *Dir) InitializedFor(dialect string) bool {
    for _, d := range d.Dialects {
        if d == dialect {
            return true
        }
    }
    return false
}

// HasMigration checks if migration already exists.
func (d *Dir) HasMigration(fileName string) bool {
    for _, fn := range d.Migrations {
        if fn == fileName {
            return true
        }
    }
    return false
}

// setErr sets error.
func (d *Dir) setErr(err error) error {
    if err != nil {
        d.Err = err
        d.Valid = false
        return err
    }
    return nil
}
