package dialect

// create migrations table.
var mysqlMigCreate = `CREATE TABLE IF NOT EXIST migrations (
  id varchar(36) NOT NULL,
  desc varchar(30) DEFAULT NULL,
  created_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB;`
