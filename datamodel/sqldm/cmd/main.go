package main

import (
	"fmt"
	"github.com/csby/tool/datamodel/sqldm/lib"
)

func main() {
	generator := &lib.Generator{Database: cfg.Database.Database()}
	err := generator.CreateEntity(cfg.Package.EntityPkg(), cfg.Package.ModelPgk())
	if err != nil {
		fmt.Println(err)
	}
}
