package mig

import (
    "regexp"
    "errors"
    "path"
    "os"
    "fmt"
    "time"
)

// Regexp pattern matching migration file name.
var migName = regexp.MustCompile(`^mig_([a-z]+)_([0-9]{19})\.go$`)

// IsMigration returns true if path is a migration file.
func IsMigration(path string) bool {
    p := migName.FindAllStringSubmatch(path, 3)
    return len(p) > 0
}

// FileDialectAndTs returns dialect and creation timestamp for given migration file.
func FileDialectAndTs(file string) (string, string, error) {
    p := migName.FindAllStringSubmatch(path.Base(file), 3)
    if len(p) == 0 && len(p[0]) == 3 {
        return "", "", errors.New("not a migration file")
    }
    return p[0][1], p[0][2], nil
}

// CreateDir creates directory if it does not exist.
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
    fi, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }
    return !fi.IsDir(), nil
}

// DirExists returns true if path points to an existing directory.
func DirExists(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        }
        return false, err
    }
    return fi.IsDir(), nil
}

// GenMigMainFileName returns main migration file name for for given dialect.
func GenMigMainFileName(dialect string) string {
    return fmt.Sprintf("mig_%s.go", dialect)
}

// GenMigFileName returns migration file name for for given dialect.
func GenMigFileName(dialect string) (int64, string) {
    ts := time.Now().UnixNano()
    return ts, fmt.Sprintf("mig_%s_%d.go", dialect, ts)
}
