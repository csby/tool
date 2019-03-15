package main

import (
	"fmt"
	"github.com/csby/database/sqldb/mssql"
	"github.com/csby/database/sqldb/mysql"
	"github.com/csby/tool/datamodel/sqldm/lib"
	"github.com/csby/tool/module"
	"os"
	"path/filepath"
)

const (
	moduleType    = "tool"
	moduleName    = "sqldmcmd"
	moduleRemark  = "数据库实体模型生成工具"
	moduleVersion = "1.0.1.0"
)

var (
	args = &Args{}
	cfg  = &Config{
		Database: ConfigDatabase{
			MySql: ConfigDatabaseMysql{
				Enable: true,
				Connection: mysql.Connection{
					Server:   "172.0.0.1",
					Port:     3306,
					Schema:   "test",
					Charset:  "utf8",
					Timeout:  10,
					User:     "root",
					Password: "",
				},
			},
			MsSql: ConfigDatabaseMssql{
				Enable: false,
				Connection: mssql.Connection{
					Server:   "127.0.0.1",
					Port:     1433,
					Schema:   "test",
					Instance: "MSSQLSERVER",
					User:     "sa",
					Password: "",
					Timeout:  10,
				},
			},
		},
		Package: ConfigPackage{
			Entity: ConfigEntity{
				Enable: true,
				Package: lib.Package{
					Name:   "entity",
					Path:   "github.com/test/data/entity",
					Folder: "src/github.com/test/data/entity",
				},
			},
			Model: ConfigEntity{
				Enable: false,
				Package: lib.Package{
					Name:   "model",
					Path:   "github.com/test/data/model",
					Folder: "src/github.com/test/data/model",
				},
			},
		},
	}
)

func init() {
	moduleArgs := &module.ArgsModule{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, args)
	rootFolder, err := filepath.Abs("cfg")
	if err != nil {
		rootFolder = filepath.Join(filepath.Dir(moduleArgs.ModuleFolder()), "cfg")
	}
	if args.help {
		args.ShowHelp(rootFolder)
		os.Exit(0)
	}

	cfgPath := args.cfg
	if cfgPath == "" {
		cfgPath = filepath.Join(rootFolder, "sqldm.json")
	} else if !filepath.IsAbs(cfgPath) {
		cfgPath = filepath.Join(rootFolder, cfgPath)
	}
	fmt.Println("cfg:", cfgPath)

	_, err = os.Stat(cfgPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(cfgPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}

	if !filepath.IsAbs(cfg.Package.Entity.Folder) {
		cfg.Package.Entity.Folder, _ = filepath.Abs(cfg.Package.Entity.Folder)
	}
	if !filepath.IsAbs(cfg.Package.Model.Folder) {
		cfg.Package.Model.Folder, _ = filepath.Abs(cfg.Package.Model.Folder)
	}
}
