package main

import (
	"encoding/json"
	"github.com/csby/database/sqldb"
	"github.com/csby/database/sqldb/mssql"
	"github.com/csby/database/sqldb/mysql"
	"github.com/csby/database/sqldb/oracle"
)

type ConfigDatabase struct {
	Enable     bool        `json:"enable"`
	Type       string      `json:"type" note:"mysql, mssql, oracle"`
	Connection interface{} `json:"connection"`
}

func (s *ConfigDatabase) Database() sqldb.SqlDatabase {
	if !s.Enable {
		return nil
	}

	if s.Type == "mysql" {
		conn := &mysql.Connection{}
		data, err := json.Marshal(s.Connection)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(data, conn)
		if err != nil {
			return nil
		}
		return mysql.NewDatabase(conn)
	} else if s.Type == "mssql" {
		conn := &mssql.Connection{}
		data, err := json.Marshal(s.Connection)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(data, conn)
		if err != nil {
			return nil
		}
		return mssql.NewDatabase(conn)
	} else if s.Type == "oracle" {
		conn := &oracle.Connection{}
		data, err := json.Marshal(s.Connection)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(data, conn)
		if err != nil {
			return nil
		}
		return oracle.NewDatabase(conn)
	} else {
		return nil
	}

	return nil
}
