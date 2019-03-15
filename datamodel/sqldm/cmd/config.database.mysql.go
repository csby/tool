package main

import "github.com/csby/database/sqldb/mysql"

type ConfigDatabaseMysql struct {
	Enable     bool             `json:"enable"`
	Connection mysql.Connection `json:"connection"`
}
