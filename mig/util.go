package mig

import (
    "regexp"
    "errors"
    "path"
    "os"
    "fmt"
    "time"
)

var migName = regexp.MustCompile(`^mig_([a-z]+)_([0-9]{19})\.go$`)

// IsMigration is file a migration file.
func IsMigration(name string) bool {
    p := migName.FindAllStringSubmatch(name, 3)
    return len(p) > 0
}

// DescMigration takes migration file name and returns dialect
// and creation timestamp.
func DescMigration(name string) (string, string, error) {
    p := migName.FindAllStringSubmatch(path.Base(name), 3)
    if len(p) == 0 {
        return "", "", errors.New("not a migration file")
    }
    return p[0][1], p[0][2], nil
}

// DirCreate creates directory if it does not exist.
func DirCreate(path string) error {
    fi, err := exists(path)
    if err != nil {
        return err
    }
    if fi == nil {
        return os.MkdirAll(path, 0777)
    }
    if !fi.IsDir() {
        return fmt.Errorf("%s already exists and it's not a directory", path)
    }
    return nil
}

// FileExists returns true if path points to an existing file.
func FileExists(path string) (bool, error) {
    fi, err := exists(path)
    if err != nil {
        return false, err
    }
    if fi != nil {
        return !fi.IsDir(), nil
    }
    return false, nil
}

// DirExists returns true if path points to an existing directory.
func DirExists(path string) (bool, error) {
    fi, err := exists(path)
    if err != nil {
        return false, err
    }
    if fi != nil {
        return fi.IsDir(), nil
    }
    return false, nil
}

func GenMigMainFileName(dialect string) string {
    return fmt.Sprintf("mig_%s.go", dialect)
}

func GenMigFileName(dialect string) (int64, string) {
    ts := time.Now().UnixNano()
    return ts, fmt.Sprintf("mig_%s_%d.go", dialect, ts)
}

// exists returns os.FileInfo if file or directory exists.
func exists(path string) (os.FileInfo, error) {
    fi, err := os.Stat(path)
    if os.IsNotExist(err) {
        return nil, nil
    }
    return fi, err
}
