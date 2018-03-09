package mig

import (
    "io/ioutil"
    "path"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// yamlConfig represents YAML configuration file.
type yamlConfig struct {
    dir      string
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
    cfg.dir = path.Dir(file)
    for n, db := range cfg.YTargets {
        db.name = n
        db.dir = path.Join(cfg.dir, "migrations", n)
    }
    return cfg, nil
}

func (c *yamlConfig) MigDir() string {
    return c.dir
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
