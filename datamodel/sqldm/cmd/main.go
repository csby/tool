package main

import (
	"fmt"
	"github.com/csby/tool/datamodel/sqldm/lib"
)

func main() {
	count := len(cfg.Items)
	fmt.Println("count:", count)
	for index := 0; index < count; index++ {
		item := cfg.Items[index]
		if item == nil {
			continue
		}

		generator := &lib.Generator{Database: item.Database.Database()}
		err := generator.CreateEntity(item.Package.EntityPkg(), item.Package.ModelPgk())
		if err != nil {
			fmt.Println(err)
		}
	}

}
