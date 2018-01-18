package mig

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

// A dir represents state of migrations directory
// and provides file manipulation methods.
type dir struct {
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
    migrations map[string]*File
    // Set to true when rootDir structure is ready to use by mig.
    valid bool
    // The error which stopped the check.
    err error
}

// newDir creates new dir instance and collects populates the state of the
// migration directory.
func newDir(root string) *dir {
    s := &dir{
        migDirEmpty: true,
        migrations:  make(map[string]*File),
    }
    s.rootDir, s.err = filepath.Abs(root)
    if s.err != nil {
        return s
    }
    s.collect()
    return s
}

// ready returns true if migration directory is ready to use by mig.
func (d *dir) ready() bool {
    return d.err == nil && d.valid
}

// init initialize migrations directory for given dialect.
func (d *dir) init(dialect string) error {
    if IsSupDialect(dialect) == false {
        return d.failure(fmt.Errorf("unsupported dialect: %d", dialect))
    }
    if d.isInitFor(dialect) == true {
        return nil
    }
    if d.migDirExist == false {
        if err := d.createMigDir(); err != nil {
            return d.failure(err)
        }
    }
    file, err := NewMigFile(d.migDir, dialect, kindDial)
    if err != nil {
        return err
    }
    return file.Save()
}

// newMigration creates new migration file for given dialect.
func (d *dir) newMigration(dialect string) (string, error) {
    if d.isInitFor(dialect) == false {
        return "", d.failure(fmt.Errorf("migrations not initialized for %s", dialect))
    }

    if d.ready() == false {
        return "", d.failure(fmt.Errorf("migrations directory not ready"))
    }

    if IsSupDialect(dialect) == false {
        return "", d.failure(fmt.Errorf("unsupported dialect: %d", dialect))
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
func (d *dir) createMigDir() error {
    if err := os.MkdirAll(d.migDir, 0777); err != nil {
        return d.failure(err)
    }
    d.collect()
    return nil
}

// collect collects all the information and sets all fields accordingly.
func (d *dir) collect() error {
    var err error
    if d.rootDirExist, err = IsDir(d.rootDir); err != nil {
        return d.failure(err)
    }

    if d.rootDirExist == false {
        return d.failure(fmt.Errorf("%d must exist", d.rootDir))
    }

    // Check migration directory exists in root and it'd actually a directory.
    var ex bool
    d.migDir = filepath.Join(d.rootDir, "migration")
    if ex, err = FileExists(d.migDir); err != nil {
        return d.failure(err)
    }
    if d.migDirExist, err = IsDir(d.migDir); err != nil {
        return d.failure(err)
    }
    if ex && !d.migDirExist {
        return d.failure(fmt.Errorf("%d is not a directory", d.migDir))
    }

    if d.migDirExist == false {
        d.valid = true
        return nil
    }

    var fs []os.FileInfo
    if fs, err = ioutil.ReadDir(d.migDir); err != nil {
        return d.failure(err)
    }
    for _, fi := range fs {
        switch {
        case IsDialectFile(fi.Name()):
            dialect, _ := DecodeDialectFile(fi.Name())
            if !IsSupDialect(dialect) {
                return d.failure(fmt.Errorf("unsupported dialect %d in "+
                    "migration file %d", dialect, fi.Name()))
            }
            d.dialects = append(d.dialects, dialect)
        case IsMigFile(fi.Name()):
            file, err := NewFileFromPath(fi.Name())
            if err != nil {
                return err
            }
            d.migrations[file.descriptor()] = file
        default:
            return d.failure(fmt.Errorf("unexpected file %d", fi.Name()))
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

// isInitFor returns true if migrations have been
// initialized for given dialect.
func (d *dir) isInitFor(dialect string) bool {
    for _, d := range d.dialects {
        if d == dialect {
            return true
        }
    }
    return false
}

// hasMigration checks if migration already exists.
func (d *dir) hasMigration(id int64, dialect string) bool {
    _, ok := d.migrations[Desc(id, dialect)]
    return ok
}

// failure sets err and makes the directory as invalid.
func (d *dir) failure(err error) error {
    if err == nil {
        return nil
    }
    d.err = err
    d.valid = false
    return err
}
