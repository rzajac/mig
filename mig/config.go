package mig

import (
    "path"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
)

// List of known database dialects.
const DialectMySQL = "mysql"

// config represents YAML configuration file.
type config struct {
    fs  afero.Fs
    Dir string `yaml:"dir"` // Absolute migrations path.
    Targets map[string]*struct {
        Dialect string `yaml:"dialect"` // Database dialect.
        Dsn     string `yaml:"dsn"`     // Database connection string.
    } `yaml:"targets"`      // List of defined targets.
}

// NewConfig loads YAML configuration.
func NewConfig(fs afero.Fs, configPath string) (*config, error) {
    content, err := afero.ReadFile(fs, configPath)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    cfg := &config{fs: fs}
    if err = yaml.UnmarshalStrict(content, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    // In config we get only directory name.
    // Here we set it as absolute and relative to configuration path.
    cfg.Dir = path.Join(path.Dir(configPath), cfg.Dir)
    return cfg, nil
}

func (c *config) MigDir() string {
    return c.Dir
}

func (c *config) Target(trgName string) (Target, error) {
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
    return NewTarget(c.fs, path.Join(c.Dir, trgName), drv, GetMigrations(trgName))
}
