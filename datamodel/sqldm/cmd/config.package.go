package main

import "github.com/csby/tool/datamodel/sqldm/lib"

type ConfigPackage struct {
	Entity ConfigEntity `json:"entity"`
	Model  ConfigEntity `json:"model"`
}

func (s *ConfigPackage) EntityPkg() *lib.Package {
	if s.Entity.Enable {
		return &s.Entity.Package
	}

	return nil
}

func (s *ConfigPackage) ModelPgk() *lib.Package {
	if s.Model.Enable {
		return &s.Model.Package
	}

	return nil
}

type ConfigEntity struct {
	Enable bool `json:"enable"`
	lib.Package
}
