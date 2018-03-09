package mig

import "github.com/pkg/errors"

const DialectMySQL = "mysql"

// DriverProvider provides database drivers based on target names.
type DriverProvider struct {
    cfg Config
}

// NewDriverProvider returns new database driver provider.
func NewDriverProvider(config Config) *DriverProvider {
    return &DriverProvider{
        cfg: config,
    }
}

// Driver returns database driver for given target name.
func (dp *DriverProvider) Driver(targetName string) (Driver, error) {
    c, err := dp.cfg.Target(targetName)
    if err != nil {
        return nil, err
    }
    switch c.Dialect() {
    case DialectMySQL:
        return newMYSQLDriver(c), nil
    default:
        return nil, errors.Errorf("unknown dialect: %s", c.Dialect())
    }
}
