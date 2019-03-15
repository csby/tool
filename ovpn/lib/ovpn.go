package lib

import (
	"fmt"
	"github.com/csby/security/certificate"
	"os"
	"path/filepath"
)

type OVpn struct {
}

func (s *OVpn) CreateClientConf(cfg *Config) (string, error) {
	caCrt := &certificate.Crt{}
	err := caCrt.FromData([]byte(cfg.Ca.Cert))
	if err != nil {
		return "", fmt.Errorf("加载CA证书失败: %v", err)
	}

	caPrivate := &certificate.RSAPrivate{}
	err = caPrivate.FromData([]byte(cfg.Ca.Key), cfg.Ca.KeyPassword)
	if err != nil {
		return "", fmt.Errorf("加载CA私钥失败: %v", err)
	}

	private := &certificate.RSAPrivate{Format: "pkcs8"}
	err = private.Create(2048)
	if err != nil {
		return "", fmt.Errorf("创建私钥失败: %v", err)
	}
	public, err := private.Public()
	if err != nil {
		return "", fmt.Errorf("获取公钥失败: %v", err)
	}

	crtTemplate := &certificate.CrtTemplate{
		Organization:       "client",
		OrganizationalUnit: cfg.Client.OU,
		Locality:           cfg.Client.Locality,
		Province:           cfg.Client.Province,
		StreetAddress:      cfg.Client.StreetAddress,
		CommonName:         cfg.Client.CN,
	}
	template, err := crtTemplate.Template()
	if err != nil {
		return "", fmt.Errorf("获取证书失败: %v", err)
	}
	crt := &certificate.Crt{}
	err = crt.Create(template, caCrt.Certificate(), public, caPrivate)
	if err != nil {
		return "", fmt.Errorf("创建证书失败: %v", err)
	}

	cert, err := crt.ToMemory()
	if err != nil {
		return "", err
	}
	key, err := private.ToMemory("")
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(cfg.Client.Folder, 0777)
	if err != nil {
		return "", err
	}
	confFilePath := filepath.Join(cfg.Client.Folder, fmt.Sprintf("%s.ovpn", cfg.Client.OU))
	file, err := os.Create(confFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, cfg.Template, cfg.Ca.Cert, cfg.Ta, string(cert), string(key))
	if err != nil {
		return "", err
	}

	return confFilePath, nil
}
