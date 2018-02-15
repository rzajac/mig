package mig

type Cfg struct {
    Dir string  `yaml:"dir"`
    DBs []DBCfg `yaml:"dbs"`
}

type DBCfg struct {
    Name    string `yaml:"name"`
    Dialect string `yaml:"dialect"`
}
