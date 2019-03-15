package config

func (s *Config) binaryFilesForCloud() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_vsgw_cloud_service": "vsgw.cloud",
	}
}

func (s *Config) binaryFilesForGateway() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_vsgw_gateway_service": "vsgw",
	}
}

func (s *Config) binaryFilesForCrtMgr() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_tool_certmaker": "certmaker",
	}
}

func (s *Config) enableAppForCrtMgr() bool {
	return false
}

func (s *Config) binaryFilesForSlqDM() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_tool_datamodel_sqldm_cmd": "sqldm",
	}
}
