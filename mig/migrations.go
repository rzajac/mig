package mig

import (
    "fmt"
    "sort"

    "github.com/pkg/errors"
)

// Registered migrations.
var migrations = make(map[string][]Migration)

// Register registers Migration.
func Register(target string, migration Migration) {
    migrations[target] = append(migrations[target], migration)
}

// validateMigs validates migrations list for given target.
func validateMigs(target string) error {
    // No migrations no possibility for error.
    if len(migrations[target]) == 0 {
        return nil
    }
    prev := migrations[target][0].AppliedAt().IsZero()
    for _, mgr := range migrations[target] {
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

// sortMigs sots all registered migrations for all targets.
func sortMigs() {
    for _, mgr := range migrations {
        sort.Sort(migSort(mgr))
    }
}

// applyAllMigs apply all outstanding migrations for given target.
func applyAllMigs(target string) error {
    for _, mgr := range migrations[target] {
        if !mgr.AppliedAt().IsZero() {
            continue
        }
        if err := mgr.Apply(); err != nil {
            return err
        }
    }
    return nil
}

// listMigs lists migrations for given target.
func listMigs(target string) {
    for _, m := range migrations[target] {
        fmt.Println(target, m.Version())
    }
}
