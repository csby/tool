package lib

import (
	"fmt"
	"github.com/csby/database/sqldb"
	"github.com/csby/database/sqldb/mssql"
	"github.com/csby/database/sqldb/mysql"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestGenerator_CreateEntity(t *testing.T) {
	rootFolder := pkgRootFolder()
	t.Log("root folder:", rootFolder)

	entityPkg := Package{
		Name: "db",
	}
	rootPath := fmt.Sprintf("%s/pgk", reflect.TypeOf(entityPkg).PkgPath())
	t.Log("root path:", rootPath)

	entityPkg.Folder = filepath.Join(rootFolder, entityPkg.Name)
	entityPkg.Path = fmt.Sprintf("%s/%s", rootPath, entityPkg.Name)
	t.Log("entity folder:", entityPkg.Folder)
	t.Log("entity path:", entityPkg.Path)

	modelPkg := Package{
		Name: "model",
	}
	modelPkg.Folder = filepath.Join(rootFolder, modelPkg.Name)
	modelPkg.Path = fmt.Sprintf("%s/%s", rootPath, modelPkg.Name)
	t.Log("model folder:", modelPkg.Folder)
	t.Log("model path:", modelPkg.Path)

	db := SqlServer()
	db = MySql()

	//createEntity(t, db, &entityPkg, nil)
	//createEntity(t, db, nil, &modelPkg)
	createEntity(t, db, &entityPkg, &modelPkg)
}

func TestGenerator_End(t *testing.T) {
	rootFolder := pkgRootFolder()
	os.RemoveAll(rootFolder)
}

func createEntity(t *testing.T, db sqldb.SqlDatabase, entity, model *Package) {
	generator := &Generator{Database: db}
	err := generator.CreateEntity(entity, model)
	if err != nil {
		t.Fatal(err)
	}
}

func pkgRootFolder() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(file), "pgk")
}

func SqlServer() sqldb.SqlDatabase {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) < 2 {
		paths = strings.Split(goPath, ":")
	}
	if len(paths) > 1 {
		goPath = paths[0]
	}
	cfgPath := filepath.Join(goPath, "tmp", "cfg", "tool_datamodel_sqldm_mssql_test.json")
	cfg := &mssql.Connection{
		Server:   "127.0.0.1",
		Port:     1433,
		Schema:   "test",
		Instance: "MSSQLSERVER",
		User:     "sa",
		Password: "",
		Timeout:  10,
	}
	_, err := os.Stat(cfgPath)
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

	return mssql.NewDatabase(cfg)
}

func MySql() sqldb.SqlDatabase {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) < 2 {
		paths = strings.Split(goPath, ":")
	}
	if len(paths) > 1 {
		goPath = paths[0]
	}
	cfgPath := filepath.Join(goPath, "tmp", "cfg", "tool_datamodel_sqldm_mysql_test.json")
	cfg := &mysql.Connection{
		Server:   "172.0.0.1",
		Port:     3306,
		Schema:   "mysql",
		Charset:  "utf8",
		Timeout:  10,
		User:     "root",
		Password: "",
	}
	_, err := os.Stat(cfgPath)
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

	return mysql.NewDatabase(cfg)
}
