package mig

import (
    "path"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
)

// List of known database dialects.
const DialectMySQL = "mysql"

// File system abstraction.
var fs = &afero.Afero{afero.NewOsFs()}

// mig represents YAML configuration file.
type mig struct {
    Dir string `yaml:"dir"` // Absolute migrations path.
    Targets map[string]*struct {
        Dialect string `yaml:"dialect"` // Database dialect.
        Dsn     string `yaml:"dsn"`     // Database connection string.
    } `yaml:"targets"`      // List of defined targets.
}

// NewMig loads YAML configuration.
func NewMig(configPath string) (*mig, error) {
    content, err := afero.ReadFile(fs, configPath)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    cfg := &mig{}
    if err = yaml.UnmarshalStrict(content, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    // In config file we only have directory name.
    // This builds absolute path to migrations directory
    // based on configuration file path.
    cfg.Dir = path.Join(path.Dir(configPath), cfg.Dir)
    return cfg, nil
}

func (c *mig) MigDir() string {
    return c.Dir
}

func (c *mig) Target(trgName string) (Target, error) {
    trgCfg, ok := c.Targets[trgName]
    if !ok {
        return nil, ErrUnknownTarget
    }

    var drv Driver
    switch trgCfg.Dialect {
    case DialectMySQL:
        drv = newMYSQLDriver(trgName, trgCfg.Dsn)
    default:
        return nil, errors.Errorf("unknown dialect: %s", trgCfg.Dialect)
    }
    return NewTarget(path.Join(c.Dir, trgName), drv, GetMigrations(trgName))
}
