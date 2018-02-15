package mig

type MigConfig struct {
    Dir string               `yaml:"dir"`
    DBs map[string]*DBConfig `yaml:"db"`
}

type DBConfig struct {
    Dialect string `yaml:"dialect"`
}
