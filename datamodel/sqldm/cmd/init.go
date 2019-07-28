package main

import (
	"fmt"
	"github.com/csby/database/sqldb/mssql"
	"github.com/csby/database/sqldb/mysql"
	"github.com/csby/database/sqldb/oracle"
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
		Items: []*ConfigItem{
			{
				Database: ConfigDatabase{
					Enable: false,
					Type:   "mysql",
					Connection: &mysql.Connection{
						Host:     "172.0.0.1",
						Port:     3306,
						Schema:   "test",
						Charset:  "utf8",
						Timeout:  10,
						User:     "root",
						Password: "",
					},
				},
				Package: ConfigPackage{
					Entity: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "entity",
							Path:   "github.com/test/mysql/data/entity",
							Folder: "src/github.com/test/mysql/data/entity",
						},
					},
					Model: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "model",
							Path:   "github.com/test/mysql/data/model",
							Folder: "src/github.com/test/mysql/data/model",
						},
					},
				},
			},
			{
				Database: ConfigDatabase{
					Enable: false,
					Type:   "mssql",
					Connection: &mssql.Connection{
						Host:     "172.0.0.1",
						Port:     1433,
						Schema:   "test",
						Instance: "MSSQLSERVER",
						User:     "sa",
						Password: "",
						Timeout:  10,
					},
				},
				Package: ConfigPackage{
					Entity: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "entity",
							Path:   "github.com/test/mssql/data/entity",
							Folder: "src/github.com/test/mssql/data/entity",
						},
					},
					Model: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "model",
							Path:   "github.com/test/mssql/data/model",
							Folder: "src/github.com/test/mssql/data/model",
						},
					},
				},
			},
			{
				Database: ConfigDatabase{
					Enable: false,
					Type:   "oracle",
					Connection: &oracle.Connection{
						Host:     "172.0.0.1",
						Port:     1521,
						SID:      "orcl",
						User:     "orc",
						Password: "",
						Owners: []string{
							"LAB",
							"EXAM",
						},
					},
				},
				Package: ConfigPackage{
					Entity: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "entity",
							Path:   "github.com/test/oracle/data/entity",
							Folder: "src/github.com/test/oracle/data/entity",
						},
					},
					Model: ConfigEntity{
						Enable: false,
						Package: lib.Package{
							Name:   "model",
							Path:   "github.com/test/oracle/data/model",
							Folder: "src/github.com/test/oracle/data/model",
						},
					},
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

	count := len(cfg.Items)
	for index := 0; index < count; index++ {
		item := cfg.Items[index]
		if item == nil {
			continue
		}

		if !filepath.IsAbs(item.Package.Entity.Folder) {
			item.Package.Entity.Folder, _ = filepath.Abs(item.Package.Entity.Folder)
		}
		if !filepath.IsAbs(item.Package.Model.Folder) {
			item.Package.Model.Folder, _ = filepath.Abs(item.Package.Model.Folder)
		}

	}
}
