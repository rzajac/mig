package mig

import (
    "time"
    "strconv"
    "path"
    "io/ioutil"
    "os"
    "text/template"
    "fmt"
    "github.com/rzajac/mig/mig/dialect"
    "bytes"
    "github.com/pkg/errors"
)

const kindMigr = "migration"
const kindDial = "dialect"

// A File represents migration file.
type File struct {
    Path    string    // Absolute path.
    Dialect string    // Migration Dialect.
    Ts      time.Time // Creation time.
    Kind    string    // File kind.
    loaded  bool      // True if loaded from filesystem.
}

// NewMigFile returns new instance of File.
// The migDir must be absolute path to migration directory.
func NewMigFile(migDir, dialect, kind string) (*File, error) {
    if IsSupDialect(dialect) == false {
        return nil, fmt.Errorf("uknknown migration dialect %s", dialect)
    }
    if IsValidKind(kind) == false {
        return nil, fmt.Errorf("uknknown migration file kind %s", kind)
    }

    var ts int64
    file := &File{Dialect: dialect, Kind: kind}

    switch file.Kind {
    case kindMigr:
        file.Path, ts = GenMigFileName(dialect)
    case kindDial:
        file.Path = GenDialectFileName(dialect)
    }

    file.Path = path.Join(migDir, file.Path)
    file.Ts = time.Unix(0, ts)
    return file, nil
}

// NewFileFromPath creates File instance from path.
// Path must be absolute path to file.
func NewFileFromPath(path string) (*File, error) {
    ex, err := FileExists(path)
    switch {
    case err != nil:
        return nil, err
    case ex == false:
        return nil, errors.New("path must point to existing migration file")
    }

    f := &File{Path: path, loaded: true}
    switch {
    case IsMigFile(f.Path):
        f.Kind = kindMigr
        var tss string
        var tsi int64
        if f.Dialect, tss, err = DecodeMigFile(f.Path); err != nil {
            return nil, err
        }
        tsi, _ = strconv.ParseInt(tss, 10, 64)
        f.Ts = time.Unix(0, tsi)
    case IsDialectFile(f.Path):
        f.Kind = kindDial
        if f.Dialect, err = DecodeDialectFile(f.Path); err != nil {
            return nil, err
        }
    }

    return f, nil
}

// IsLoaded returns true if file has been loaded from filesystem.
func (f *File) IsLoaded() bool {
    return f.loaded
}

// Delete deletes the file.
func (f *File) Delete() error {
    return os.Remove(f.Path)
}

// Save saves new dialect or migration file.
func (f *File) Save() error {
    if f.loaded {
        return errors.New("overwriting migration files is forbidden")
    }
    content, err := f.genContent()
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(f.Path, content, 0666); err != nil {
        return err
    }
    return nil
}

// genContent returns dialect or migration file content.
func (f *File) genContent() ([]byte, error) {
    var tpl *template.Template
    var buf bytes.Buffer

    switch f.Dialect {
    case "mysql":
        if f.Kind == kindDial {
            tpl = dialect.MysqlDialectTpl
        } else {
            tpl = dialect.MysqlMigTpl
        }
    }

    if err := tpl.Execute(&buf, f); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

