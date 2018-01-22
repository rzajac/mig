package mig

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "text/template"
    "time"
)

// A MigFile represents migration file.
type MigFile struct {
    Path    string    // Absolute path.
    Dialect string    // Migration Dialect.
    Ts      time.Time // Creation time.
    loaded  bool      // True if loaded from filesystem.
}

// newMigFile returns new instance of MigFile.
// The migDir must be absolute path to migration directory.
func newMigFile(migDir, dialect string) *MigFile {
    var ts int64
    file := &MigFile{Dialect: dialect}
    file.Path, ts = NextMigFileName(dialect)
    file.Path = path.Join(migDir, file.Path)
    file.Ts = time.Unix(0, ts)
    return file
}

// newFileFromPath creates MigFile instance from path to existing migration file.
// Path must be absolute path to migration file.
func newFileFromPath(path string) (*MigFile, error) {
    var err error
    var tsi int64

    f := &MigFile{Path: path, loaded: true}
    f.Dialect, tsi, err = ParseMigFileName(f.Path)
    if err != nil {
        return nil, err
    }
    f.Ts = time.Unix(0, tsi)
    return f, nil
}

// IsLoaded returns true if file has been loaded from filesystem.
func (f *MigFile) IsLoaded() bool {
    return f.loaded
}

// Delete deletes the file.
func (f *MigFile) Delete() error {
    return os.Remove(f.Path)
}

// Save saves new migration file to migration directory.
func (f *MigFile) Save() error {
    if f.loaded {
        return fmt.Errorf("saving existing migration file %s is forbidden", f.Path)
    }
    content, err := f.getTemplate()
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(f.Path, content, 0666); err != nil {
        return err
    }
    return nil
}

// getTemplate returns migration file template.
func (f *MigFile) getTemplate() ([]byte, error) {
    var tpl *template.Template
    var buf bytes.Buffer

    switch f.Dialect {
    case "mysql":
        tpl = mySQLMigrationTpl
    }

    if err := tpl.Execute(&buf, f); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// descriptor returns unique migration file descriptor.
func (f *MigFile) descriptor() string {
    return MigrationDescriptor(f.Dialect, f.Ts.UnixNano())
}
