package mig

import (
    "io/ioutil"
    "path"

    "github.com/pkg/errors"
    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
)

// List of known database dialects.
const DialectMySQL = "mysql"

// config represents YAML configuration file.
type config struct {
    Dir     string             `yaml:"dir"`     // Absolute migrations path.
    Targets map[string]*target `yaml:"targets"` // List of defined targets.
}

// NewYAMLCfg loads YAML configuration.
func NewYAMLCfg(file string) (*config, error) {
    content, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    cfg := &config{}
    if err = yaml.UnmarshalStrict(content, cfg); err != nil {
        return nil, errors.WithStack(err)
    }

    filePath := path.Dir(file)
    for name, target := range cfg.Targets {
        target.name = name
        target.dir = path.Join(filePath, cfg.Dir, name)
    }
    return cfg, nil
}

func (c *config) MigDir() string {
    return c.Dir
}

func (c *config) Driver(target string) (Driver, error) {
    trg, ok := c.Targets[target]
    if !ok {
        return nil, errors.Errorf("no database target named %s", target)
    }
    switch trg.Dialect() {
    case DialectMySQL:
        drv := newMYSQLDriver(trg)
        return drv, drv.Open()
    default:
        return nil, errors.Errorf("unknown dialect: %s", trg.Dialect())
    }
}

func (c *config) TargetNames() ([]string) {
    var t []string
    for n := range c.Targets {
        t = append(t, n)
    }
    return t
}

// target represents configuration of the migration target.
type target struct {
    name      string                  // Target name.
    dir       string                  // Absolute path to migrations directory.
    DbDialect string `yaml:"dialect"` // Database dialect.
    DbDsn     string `yaml:"dsn"`     // Database connection string.
}

func (t *target) Dialect() string {
    return t.DbDialect
}

func (t *target) Dsn() string {
    return t.DbDsn
}

func (t *target) Name() string {
    return t.name
}

func (t *target) MigDir() string {
    return t.dir
}
