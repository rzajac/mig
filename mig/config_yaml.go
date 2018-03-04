package mig

import (
    "io/ioutil"
    "path"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// yamlConfig represents YAML configuration file.
type yamlConfig struct {
    migDir   string
    YTargets map[string]*target `yaml:"targets"`
}

// NewYAMLCfg loads configuration from YAML file and returns yamlConfig.
func NewYAMLCfg(file string) (*yamlConfig, error) {
    var err error
    var data []byte
    cfg := &yamlConfig{}
    if data, err = ioutil.ReadFile(file); err != nil {
        return nil, errors.WithStack(err)
    }
    if err = yaml.UnmarshalStrict(data, cfg); err != nil {
        return nil, errors.WithStack(err)
    }
    cfg.migDir = path.Dir(file)
    for n, db := range cfg.YTargets {
        db.name = n
        db.migDir = path.Join(cfg.migDir, "migrations", n)
    }
    return cfg, nil
}

func (c *yamlConfig) MigDir() string {
    return c.migDir
}

func (c *yamlConfig) Target(name string) (Target, error) {
    db, ok := c.YTargets[name]
    if !ok {
        return nil, errors.Errorf("no database config named %s", name)
    }
    return db, nil
}

func (c *yamlConfig) Targets() ([]string) {
    var t []string
    for n := range c.YTargets {
        t = append(t, n)
    }
    return t
}
