package mig

import (
    "os"
    "strconv"
    "text/template"
    "time"

    "github.com/pkg/errors"
)

// checkCreateDir creates directory if doesn't exist.
func checkCreateDir(path string) error {
    ok, err := isDir(path)
    if err != nil {
        return err
    }
    if !ok {
        if err := os.MkdirAll(path, 0777); err != nil {
            return errors.WithStack(err)
        }
    }
    return nil
}

// isDir returns true if path points to an existing directory.
func isDir(path string) (bool, error) {
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

// fileExists returns true if path points to an existing file.
func fileExists(path string) (bool, error) {
    fi, err := os.Stat(path)
    switch {
    case os.IsNotExist(err):
        return false, nil
    case err != nil:
        return false, errors.WithStack(err)
    default:
        return !fi.IsDir(), nil
    }
}
