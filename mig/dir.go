package mig

import (
    "bytes"
    "io/ioutil"
    "os"
    "path/filepath"
    "text/template"
)

// A dir represents directory where migrations are stored and provides
// methods which manipulate on it.
type dir struct {
    // The absolute path to directory with migration files.
    migDir string
    // The list of migrations.
    migrations map[string]*MigFile
}

// newDir returns new dir instance and loads all migrations.
// If path does not exist it will attempt to create it.
func newDir(path string) (*dir, error) {
    d := &dir{}
    d.migDir = path
    d.migrations = make(map[string]*MigFile)
    if err := os.MkdirAll(d.migDir, 0777); err != nil {
        return nil, err
    }
    return d, d.loadMigrations()
}

// init initialize migrations directory for given dialect.
func (d *dir) init(dialect string) error {
    var tpl *template.Template
    switch dialect {
    case "mysql":
        tpl = mySQLStructTpl
    }

    var buf bytes.Buffer
    if err := tpl.Execute(&buf, nil); err != nil {
        return err
    }
    path := filepath.Join(d.migDir, StructFileName(dialect))
    return ioutil.WriteFile(path, buf.Bytes(), 0666)
}

// newMigration creates new migration file for given dialect and returns
// absolute path to it.
func (d *dir) newMigration(dialect string) (string, error) {
    ok, err := d.isInitFor(dialect)
    if err != nil {
        return "", err
    }
    if !ok {
        if err := d.init(dialect); err != nil {
            return "", err
        }
    }

    file := newMigFile(d.migDir, dialect)
    if err := file.Save(); err != nil {
        return "", err
    }
    d.migrations[file.descriptor()] = file

    return file.Path, nil
}

// loadMigrations loads migration files.
func (d *dir) loadMigrations() error {
    var err error
    var fs []os.FileInfo

    if fs, err = ioutil.ReadDir(d.migDir); err != nil {
        return err
    }
    for _, fi := range fs {
        if !IsMigFile(fi.Name()) {
            continue
        }
        file, err := newFileFromPath(fi.Name())
        if err != nil {
            return err
        }
        d.migrations[file.descriptor()] = file
    }
    return nil
}

// isInitFor returns true if migrations have been
// initialized for given dialect.
func (d *dir) isInitFor(dialect string) (bool, error) {
    path := filepath.Join(d.migDir, StructFileName(dialect))
    return FileExists(path)
}

// hasMigration checks if migration already exists.
func (d *dir) hasMigration(id int64, dialect string) bool {
    _, ok := d.migrations[MigrationDescriptor(dialect, id)]
    return ok
}
