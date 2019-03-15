package main

import (
	"github.com/csby/database/sqldb"
	"github.com/csby/database/sqldb/mssql"
	"github.com/csby/database/sqldb/mysql"
)

type ConfigDatabase struct {
	MySql ConfigDatabaseMysql `json:"mysql"`
	MsSql ConfigDatabaseMssql `json:"mssql"`
}

func (s *ConfigDatabase) Database() sqldb.SqlDatabase {
	if s.MySql.Enable {
		return mysql.NewDatabase(&s.MySql.Connection)
	} else if s.MsSql.Enable {
		return mssql.NewDatabase(&s.MsSql.Connection)
	}

	return nil
}
