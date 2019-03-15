package config

func (s *Config) binaryFilesForCloud() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_vsgw_cloud_service.exe": "vsgw.cloud.exe",
	}
}

func (s *Config) binaryFilesForGateway() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_vsgw_gateway_service.exe": "vsgw.exe",
	}
}

func (s *Config) binaryFilesForCrtMgr() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_tool_certmaker.exe": "certmaker.exe",
	}
}

func (s *Config) enableAppForCrtMgr() bool {
	return true
}

func (s *Config) binaryFilesForSlqDM() map[string]string {
	return map[string]string{
		"go_build_github_com_csby_tool_datamodel_sqldm_cmd.exe": "sqldm.exe",
	}
}
