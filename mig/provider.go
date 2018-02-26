package mig

import "github.com/pkg/errors"

type provider struct {
    cfg Configurator
}

// NewProvider returns new database driver provider.
func NewProvider(config Configurator) DriverProvider {
    return &provider{config}
}

// Driver returns driver for given migration configuration name.
func (p *provider) Driver(name string) (Driver, error) {
    c, err := p.cfg.DbConfig(name)
    if err != nil {
        return nil, err
    }
    switch c.Dialect() {
    case "mysql":
        return NewMYSQLDriver(c), nil
    default:
        return nil, errors.Errorf("unknown dialect: %s", c.Dialect())
    }
}
