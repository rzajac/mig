package mig

import (
    "path"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
)

// List of known database dialects.
const DialectMySQL = "mysql"

// mig represents YAML configuration file.
type mig struct {
    fs  afero.Fs
    Dir string `yaml:"dir"` // Absolute migrations path.
    Targets map[string]*struct {
        Dialect string `yaml:"dialect"` // Database dialect.
        Dsn     string `yaml:"dsn"`     // Database connection string.
    } `yaml:"targets"`      // List of defined targets.
}

// NewMig loads YAML configuration.
func NewMig(fs afero.Fs, configPath string) (*mig, error) {
    content, err := afero.ReadFile(fs, configPath)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    cfg := &mig{fs: fs}
    if err = yaml.UnmarshalStrict(content, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    // In config file have only directory name.
    // Here we set it as absolute path.
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
    return NewTarget(c.fs, path.Join(c.Dir, trgName), drv, GetMigrations(trgName))
}
