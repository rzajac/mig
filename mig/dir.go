package mig

import (
    "path/filepath"
    "os"
    "io/ioutil"
    "fmt"
    "path"
)

// A Dir represents migrations directory.
type Dir struct {
    path string // Absolute path to migrations directory.
}

// NewDir creates new instance and ensures it is a directory.
// I will create directory if it does not exist.
func NewDir(dir string) (*Dir, error) {
    var err error
    d := &Dir{}
    if d.path, err = filepath.Abs(dir); err != nil {
        return nil, err
    }
    d.path = filepath.Join(d.path, "migration")
    return d, nil
}

func (d *Dir) Initialize(dialect string) error {
    fmt.Printf("initializing migrations in %s\n", d.path)
    ex, err := d.Exists()
    if err != nil {
        return err
    }
    if ex {
        c, err := d.FileCount()
        if err != nil {
            return err
        }
        if c > 0 {
            return fmt.Errorf("cannot initialize in non empty directory %s", d.path)
        }
    }
    if err = DirCreate(d.path); err != nil {
        return err
    }
    switch dialect {
    case "mysql":
        fn := path.Join(d.path, GenMigMainFileName(dialect))
        return ioutil.WriteFile(fn, []byte(migTplMainMysql), 0666)
    default:
        panic("unknown dialect " + dialect)
    }
    return nil
}

func (d *Dir) NewMigration(dialect string) error {
    fmt.Printf("creating migration file for %s in %s\n", dialect, d.path)
    ex, err := d.Exists()
    if err != nil {
        return err
    }
    if !ex {
        return fmt.Errorf("migration directory %s does not exist\n", d.path)
    }
    switch dialect {
    case "mysql":
        ts, fn := GenMigFileName(dialect)
        fn = path.Join(d.path, fn)
        data := fmt.Sprintf(migTplMigTpl, ts, ts)
        return ioutil.WriteFile(fn, []byte(data), 0666)
    default:
        panic("unknown dialect " + dialect)
    }
    return nil
}

func (d *Dir) Exists() (bool, error) {
    return DirExists(d.path)
}

func (d *Dir) Create() error {
    return DirCreate(d.path)
}

func (d *Dir) MigCount() (int, error) {
    fs, err := d.files()
    if err != nil {
        return 0, err
    }
    fc := 0
    for _, f := range fs {
        _, err := FileFromPath(f.Name())
        if err != nil {
            fc++
        }
    }
    return fc, nil
}

// FileCount returns number of files in the directory.
func (d *Dir) FileCount() (int, error) {
    fs, err := d.files()
    if err != nil {
        return 0, err
    }
    return len(fs), nil
}

// files returns all the migration files in the directory.
func (d *Dir) files() ([]os.FileInfo, error) {
    return ioutil.ReadDir(d.path)
}
