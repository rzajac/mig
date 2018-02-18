package mig

import (
    "go/build"
    "io/ioutil"
    "os"
    "path"

    "gopkg.in/yaml.v2"
)

// Config represents Mig configuration.
type Config struct {
    Path      string                                // Absolute path to configuration file.
    Package   string               `yaml:"package"` // Migrations root directory.
    Databases map[string]*DBConfig `yaml:"db"`      // List of database configurations.
}

// DBConfig represents database configuration.
type DBConfig struct {
    Config  Config
    Dialect string `yaml:"dialect"`
    Host    string `yaml:"host"`
    User    string `yaml:"user"`
    Pass    string `yaml:"pass"`
    Name    string `yaml:"name"`
}

// GetDBConfig returns database configuration by name.
func (cfg *Config) GetDBConfig(name string) *DBConfig {
    if db, ok := cfg.Databases[name]; ok {
        return db
    }
    return nil
}

// validate validates mig configuration.
func (cfg *Config) validate() error {
    return nil
}

// MigDir returns absolute path to the directory for migrations for given
// database name.
func (cfg *DBConfig) MigDir(name string) (string, error) {
    return path.Join(cfg.Path, db.Name, db.Dialect), nil
}

// LoadConfig loads mig configuration file located at cfgPath.
func LoadConfig(cfgPath string) (*Config, error) {
    file, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        return nil, err
    }
    cfg := &Config{}
    if err := yaml.UnmarshalStrict(file, cfg); err != nil {
        return nil, err
    }
    gp := os.Getenv("GOPATH")
    if gp == "" {
        gp = build.Default.GOPATH
    }
    cfg.Path = path.Join(gp, "src", cfg.Package)
    if err := cfg.validate(); err != nil {
        return nil, err
    }
    return cfg, nil
}
