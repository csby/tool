package main

type ConfigCa struct {
	CrtFile     string `json:"crtFile" note:"证书文件路径"`
	KeyFile     string `json:"keyFile" note:"私钥文件路径"`
	KeyPassword string `json:"keyPassword" note:"私钥密码"`
}

func (s *ConfigCa) Password() *ConfigPassword {
	return &ConfigPassword{value: &s.KeyPassword}
}
