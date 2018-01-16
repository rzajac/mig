package mig

import (
    "io/ioutil"
    "os"

    "github.com/pkg/errors"
)

// ensureDir ensures path is a directory.
// If it does not exist it will create one.
func ensureDir(path string) error {
    fi, err := os.Stat(path)
    if os.IsNotExist(err) {
        return os.Mkdir(path, 0777)
    }
    if err != nil {
        return err
    }
    switch mode := fi.Mode(); {
    case mode.IsDir():
        return nil
    default:
        return errors.New(path + " not a directory")
    }
    return nil
}

// migCount returns number of files in the directory.
func migCount(path string) (int, error) {
    fs, err := ioutil.ReadDir(path)
    if err != nil {
        return 0, err
    }
    return len(fs), nil
}
