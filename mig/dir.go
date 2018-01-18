package mig

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

// A Dir represents state of migrations directory
// and provides file manipulation methods.
type Dir struct {
    // Absolute path root directory.
    // The root directory is defined as
    // the one containing "migration" directory.
    rootDir string
    // The path rootDir points to existing directory.
    rootDirExist bool
    // The absolute path to migration directory.
    // migDir is always a subdirectory of rootDir.
    migDir string
    // The path migDir points to existing directory.
    migDirExist bool
    // Set to true if migDir is empty.
    migDirEmpty bool
    // Set to true if migration directory exists and is
    // initialized with at least one Dialect.
    initialized bool
    // List of initialized dialects.
    dialects []string
    // The list of migrations.
    migrations []string
    // Set to true when rootDir structure is ready to use by mig.
    valid bool
    // The error which stopped the check.
    Error error
}

// NewDir creates new Dir instance and collects populates the state of the
// migration directory.
func NewDir(root string) *Dir {
    s := &Dir{migDirEmpty: true}
    s.rootDir, s.Error = filepath.Abs(root)
    if s.Error != nil {
        return s
    }
    s.collect()
    return s
}

// Ready returns true if migration directory is ready to use by mig.
func (d *Dir) Ready() bool {
    return d.Error == nil && d.valid
}

// Init initialize migrations directory for given dialect.
func (d *Dir) Init(dialect string) error {
    if IsSupDialect(dialect) == false {
        return d.err(fmt.Errorf("unsupported dialect: %d", dialect))
    }
    if d.IsInitFor(dialect) == true {
        return nil
    }
    if d.migDirExist == false {
        if err := d.createMigDir(); err != nil {
            return d.err(err)
        }
    }
    file, err := NewMigFile(d.migDir, dialect, kindDial)
    if err != nil {
        return err
    }
    return file.Save()
}

// NewMigration creates new migration file for given dialect.
func (d *Dir) NewMigration(dialect string) (string, error) {
    if d.IsInitFor(dialect) == false {
        return "", d.err(fmt.Errorf("migrations not initialized for %s", dialect))
    }

    if d.Ready() == false {
        return "", d.err(fmt.Errorf("migrations directory not ready"))
    }

    if IsSupDialect(dialect) == false {
        return "", d.err(fmt.Errorf("unsupported dialect: %d", dialect))
    }

    file, err := NewMigFile(d.migDir, dialect, kindMigr)
    if err != nil {
        return "", err
    }
    if err := file.Save(); err != nil {
        return "", err
    }
    return file.Path, d.collect()
}

// createMigDir creates migration directory in rootDir.
func (d *Dir) createMigDir() error {
    if err := os.MkdirAll(d.migDir, 0777); err != nil {
        return d.err(err)
    }
    d.collect()
    return nil
}

// collect collects all the information and sets all fields accordingly.
func (d *Dir) collect() error {
    var err error
    if d.rootDirExist, err = IsDir(d.rootDir); err != nil {
        return d.err(err)
    }

    if d.rootDirExist == false {
        return d.err(fmt.Errorf("%d must exist", d.rootDir))
    }

    // Check migration directory exists in root and it'd actually a directory.
    var ex bool
    d.migDir = filepath.Join(d.rootDir, "migration")
    if ex, err = FileExists(d.migDir); err != nil {
        return d.err(err)
    }
    if d.migDirExist, err = IsDir(d.migDir); err != nil {
        return d.err(err)
    }
    if ex && !d.migDirExist {
        return d.err(fmt.Errorf("%d is not a directory", d.migDir))
    }

    if d.migDirExist == false {
        d.valid = true
        return nil
    }

    var fs []os.FileInfo
    if fs, err = ioutil.ReadDir(d.migDir); err != nil {
        return d.err(err)
    }
    for _, fi := range fs {
        switch {
        case IsDialectFile(fi.Name()):
            dialect, _ := DecodeDialectFile(fi.Name())
            if !IsSupDialect(dialect) {
                return d.err(fmt.Errorf("unsupported dialect %d in "+
                    "migration file %d", dialect, fi.Name()))
            }
            d.dialects = append(d.dialects, dialect)
        case IsMigFile(fi.Name()):
            d.migrations = append(d.migrations, fi.Name())
        default:
            return d.err(fmt.Errorf("unexpected file %d", fi.Name()))
        }
    }

    if len(d.dialects) > 0 {
        d.initialized = true
    }

    if len(d.migrations) > 0 {
        d.migDirEmpty = false
    }

    d.valid = true
    return nil
}

// IsInitFor returns true if migrations have been
// initialized for given dialect.
func (d *Dir) IsInitFor(dialect string) bool {
    for _, d := range d.dialects {
        if d == dialect {
            return true
        }
    }
    return false
}

// HasMigration checks if migration already exists.
func (d *Dir) HasMigration(fileName string) bool {
    for _, fn := range d.migrations {
        if fn == fileName {
            return true
        }
    }
    return false
}

// err sets error and makes the directory as invalid.
func (d *Dir) err(err error) error {
    if err != nil {
        d.Error = err
        d.valid = false
        return err
    }
    return nil
}
