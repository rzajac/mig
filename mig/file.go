package mig

import (
    "fmt"
    "time"
)

// A File represents migration file.
type File struct {
    dir string
}

// NewMigrationFile returns new instance and prepares migrations directory.
// The dir must be a path to migrations directory.
// If directory does not exist it will create new one.
func NewMigrationFile(dir string) (*File, error) {
    dir, err := toAbs(dir)
    if err != nil {
        return nil, err
    }
    if err := ensureDir(dir); err != nil {
        return nil, err
    }
    return &File{dir}, nil
}

// Create creates new migration file.
func (m *File) Create() error {
    _, err := migCount(m.dir)
    if err != nil {
        return err
    }
    ts := time.Now().UnixNano()
    fmt.Printf("%d\n", ts)
    return nil
}
