package mig

import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "regexp"
    "time"
)

// Regexp pattern matching migration file name.
var migFileName = regexp.MustCompile(`^mig_([a-z]+)_([0-9]{19})\.go$`)

// Regexp pattern matching migration file name.
var migStructFileName = regexp.MustCompile(`^mig_(mysql)\.go$`)

// IsMigFile returns true if path is a migration file.
func IsMigFile(file string) bool {
    p := migFileName.FindAllStringSubmatch(path.Base(file), 3)
    return len(p) == 1 && len(p[0]) == 3
}

// IsMigStructFile returns true if path is a migration struct file.
func IsMigStructFile(file string) bool {
    p := migStructFileName.FindAllStringSubmatch(path.Base(file), 2)
    return len(p) == 1 && len(p[0]) == 2
}

// MigFileParts returns Dialect and creation timestamp for
// given migration file name.
func MigFileParts(file string) (string, string, error) {
    p := migFileName.FindAllStringSubmatch(path.Base(file), 3)
    if len(p) == 0 && len(p[0]) == 3 {
        return "", "", errors.New("not a migration file")
    }
    return p[0][1], p[0][2], nil
}

// StructFileParts returns Dialect for given struct migration file name.
func StructFileParts(file string) (string, error) {
    p := migStructFileName.FindAllStringSubmatch(path.Base(file), 2)
    if len(p) == 0 && len(p[0]) == 2 {
        return "", errors.New("not a migration struct file")
    }
    return p[0][1], nil
}

// CreateDir creates directory.
func CreateDir(path string) error {
    fi, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return os.MkdirAll(path, 0777)
        }
        return err
    }
    if !fi.IsDir() {
        return fmt.Errorf("%s already exists and it's not a directory", path)
    }
    return nil
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
        return false, err
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

// FileCount returns number of files in the directory.
func FileCount(dirname string) (int, error) {
    fs, err := ioutil.ReadDir(dirname)
    if err != nil {
        return 0, err
    }
    return len(fs), nil
}

// MigStructFileName returns main migration file name for for given Dialect.
func MigStructFileName(dialect string) string {
    return fmt.Sprintf("mig_%s.go", dialect)
}

// MigFileName returns migration file name for for given Dialect.
func MigFileName(dialect string) (int64, string) {
    ts := time.Now().UnixNano()
    return ts, fmt.Sprintf("mig_%s_%d.go", dialect, ts)
}

