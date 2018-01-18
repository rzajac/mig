package mig

import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "regexp"
    "time"
    "strconv"
)

// Regexp pattern matching migration file name.
var migFileName = regexp.MustCompile(`^mig_([a-z]+)_([0-9]{19})\.go$`)

// Regexp pattern matching dialect file.
// The dialect file defines the struct type for given dialect.
var dialectFileName = regexp.MustCompile(`^mig_(mysql)\.go$`)

// IsMigFile returns true if path is a migration file.
func IsMigFile(file string) bool {
    p := migFileName.FindAllStringSubmatch(path.Base(file), 3)
    return len(p) == 1 && len(p[0]) == 3
}

// DecodeMigFile decodes dialect and timestamp which are encoded in the
// migration file name. Returns dialect and migration id.
func DecodeMigFile(file string) (string, int64, error) {
    p := migFileName.FindAllStringSubmatch(path.Base(file), 3)
    if len(p) == 0 && len(p[0]) == 3 {
        return "", 0, fmt.Errorf("%s is not a migration file", file)
    }
    id, err := strconv.ParseInt(p[0][2], 10, 64)
    if err != nil {
        return "", 0, err
    }
    return p[0][1], id, nil
}

// IsDialectFile returns true if path is a dialect file.
func IsDialectFile(file string) bool {
    p := dialectFileName.FindAllStringSubmatch(path.Base(file), 2)
    return len(p) == 1 && len(p[0]) == 2
}

// DecodeDialectFile decodes dialect which is encoded in the
// migration dialect file name.
func DecodeDialectFile(file string) (string, error) {
    p := dialectFileName.FindAllStringSubmatch(path.Base(file), 2)
    if len(p) == 0 && len(p[0]) == 2 {
        return "", errors.New("not a migration struct file")
    }
    return p[0][1], nil
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

// IsValidKind returns true if the migration file kind is valid.
func IsValidKind(kind string) bool {
    return kind == kindDial || kind == kindMigr
}

// FileCount returns number of files in a given directory.
func FileCount(path string) (int, error) {
    fs, err := ioutil.ReadDir(path)
    if err != nil {
        return 0, err
    }
    return len(fs), nil
}

// GenDialectFileName generates dialect migration file name.
// Returned file name matches dialectFileName regexp.
func GenDialectFileName(dialect string) string {
    return fmt.Sprintf("mig_%s.go", dialect)
}

// GenMigFileName generates migration file name for for given dialect.
// Returned file name matches migFileName regexp.
func GenMigFileName(dialect string) (string, int64) {
    ts := time.Now().UnixNano()
    desc := Desc(ts, dialect)
    return fmt.Sprintf("mig_%s.go", desc), ts
}

// Desc returns unique migration job descriptor.
func Desc(id int64, dialect string) string {
    return dialect + "_" + strconv.FormatInt(id, 10)
}
