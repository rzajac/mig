package mig

import (
    "time"
    "path/filepath"
    "strconv"
)

// A File represents migration file.
type File struct {
    dir     *Dir      // Directory file is located in.
    path    string    // Absolute path.
    dialect string    // Migration dialect.
    ts      time.Time // Creation time.
}

// FileFromPath creates File instance from path.
func FileFromPath(path string) (*File, error) {
    path, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }
    dialect, tss, err := DescMigration(path)
    if err != nil {
        return nil, err
    }
    tsi, _ := strconv.ParseInt(tss, 10, 64)
    dir, err := NewDir(filepath.Dir(path))
    if err != nil {
        return nil, err
    }
    f := &File{
        dir:     dir,
        path:    path,
        dialect: dialect,
        ts:      time.Unix(0, tsi),
    }
    return f, nil
}
