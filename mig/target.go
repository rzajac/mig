package mig

import (
    "fmt"
    "path"
    "time"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
)

// target represents configuration of the migration target.
type target struct {
    name string      // Target name.
    dir  string      // Absolute path to migrations directory.
    drv  Driver      // The database driver.
    migs []Migration // Migrations.
}

// NewTarget returns new validated target instance.
func NewTarget(dir string, drv Driver, migs []Migration) (*target, error) {
    target := &target{
        name: path.Base(dir),
        dir:  dir,
        drv:  drv,
        migs: migs,
    }
    if err := drv.Merge(migs); err != nil {
        return nil, err
    }
    if err := target.validate(); err != nil {
        return nil, err
    }
    return target, nil
}

func (t *target) Name() string {
    return t.name
}

func (t *target) TargetDir() string {
    return t.dir
}

func (t *target) CreateMigration() error {
    if err := checkCreateDir(t.dir); err != nil {
        return err
    }
    version := time.Now().UnixNano()
    buf, err := t.drv.GenMigration(version)
    if err != nil {
        return err
    }
    migFile := path.Join(t.dir, fmt.Sprintf("%d.go", version))
    if err := afero.WriteFile(fs, migFile, buf, 0666); err != nil {
        return errors.WithStack(err)
    }
    if err := createMain(path.Dir(t.dir), t.name); err != nil {
        return err
    }
    return nil
}

func (t *target) Initialize() error {
    return t.drv.Initialize()
}

func (t *target) Migrate(toVersion int64) error {
    return nil
}

func (t *target) Status() []Status {
    stats := make([]Status, 0)
    for _, m := range t.migs {
        stats = append(stats, m)
    }
    return stats
}

// validate validates migrations for target.
func (t *target) validate() error {
    if len(t.migs) == 0 {
        return nil
    }
    prev := t.migs[0].AppliedAt().IsZero()
    for _, mgr := range t.migs {
        curr := mgr.AppliedAt().IsZero()
        switch {
        case prev == false && curr == true:
            return errors.New("migrations are not continuous")
        default:
            prev = curr
        }
    }
    return nil
}
