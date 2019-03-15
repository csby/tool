package main

import "github.com/csby/database/sqldb/mssql"

type ConfigDatabaseMssql struct {
	Enable     bool             `json:"enable"`
	Connection mssql.Connection `json:"connection"`
}
