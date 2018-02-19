package mig

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "regexp"
    "strconv"
    "text/template"
    "time"

    "github.com/pkg/errors"
)

// Regexp pattern matching migration file name.
var migrationFileName = regexp.MustCompile(`^mig_([a-z]+)_([0-9]{19})\.go$`)

// IsMigFile returns true if path is a migration file.
func IsMigFile(file string) bool {
    p := migrationFileName.FindAllStringSubmatch(path.Base(file), 3)
    return len(p) == 1 && len(p[0]) == 3
}

// ParseMigFileName returns dialect and creation timestamp
// which are part of migration file name.
func ParseMigFileName(file string) (string, int64, error) {
    p := migrationFileName.FindAllStringSubmatch(path.Base(file), 3)
    if len(p) == 0 && len(p[0]) == 3 {
        return "", 0, fmt.Errorf("%s is not a migration file", file)
    }
    id, err := strconv.ParseInt(p[0][2], 10, 64)
    if err != nil {
        return "", 0, err
    }
    return p[0][1], id, nil
}

// FileExists returns true if path points to an existing file.
func FileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    switch {
    case os.IsNotExist(err):
        return false, nil
    case err != nil:
        return false, err
    default:
        return true, nil
    }
}

// DirExists returns true if path points to an existing directory.
func IsDir(path string) (bool, error) {
    fi, err := os.Stat(path)
    switch {
    case os.IsNotExist(err):
        return false, nil
    case err != nil:
        return false, errors.WithStack(err)
    default:
        return fi.IsDir(), nil
    }
}

// IsSupDialect returns true if Dialect is on the list of supported dialects.
func IsSupDialect(dialect string) bool {
    for _, d := range dialects {
        if d == dialect {
            return true
        }
    }
    return false
}

// StructFileName returns struct file name for given dialect.
// Returned file name matches structFileName regexp.
func StructFileName(dialect string) string {
    return fmt.Sprintf("mig_%s.go", dialect)
}

// NextMigFileName returns unique migration file name.
// Returned file name matches migrationFileName regexp.
func NextMigFileName(dialect string) (string, int64) {
    ts := time.Now().UnixNano()
    desc := MigrationDescriptor(dialect, ts)
    return fmt.Sprintf("mig_%s.go", desc), ts
}

// MigrationDescriptor returns unique migration descriptor.
func MigrationDescriptor(dialect string, id int64) string {
    return dialect + "_" + strconv.FormatInt(id, 10)
}

// init initialize migrations directory for given dialect.
func CreateBaseMig(dst, dialect string) error {
    var tpl *template.Template
    switch dialect {
    case "mysql":
        tpl = mySQLStructTpl
    }
    var data struct {
        Pkg string
    }
    data.Pkg = path.Base(dst)

    var buf bytes.Buffer
    if err := tpl.Execute(&buf, data); err != nil {
        return errors.WithStack(err)
    }
    dst = filepath.Join(dst, StructFileName(dialect))
    if err := ioutil.WriteFile(dst, buf.Bytes(), 0666); err != nil {
        return errors.WithStack(err)
    }
    return nil
}
