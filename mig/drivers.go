package mig

import "github.com/pkg/errors"

// List of registered drivers.
var drivers = make(map[string]DriverMaker)

// RegisterDriver register driver.
func RegisterDriver(dialect string, maker DriverMaker) {
    drivers[dialect] = maker
}

// GetDriver returns registered driver.
func GetDriver(dialect string) (DriverMaker, error) {
    maker, ok := drivers[dialect]
    if !ok {
        return nil, errors.Errorf("unknown dialect %s", dialect)
    }
    return maker, nil
}
