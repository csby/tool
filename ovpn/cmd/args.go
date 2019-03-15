package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Args struct {
	cfg          string
	caCrt        string
	caKey        string
	caKeyPwd     string
	ta           string
	template     string
	OutputFolder string
	help         bool
}

func (s *Args) Parse(key, value string) {
	if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.help = true
	} else if key == strings.ToLower("-cfg") {
		s.cfg = value
	} else if key == strings.ToLower("-ca.crt") {
		s.caCrt = value
	} else if key == strings.ToLower("-ca.key") {
		s.caKey = value
	} else if key == strings.ToLower("-ca.key.pwd") {
		s.caKeyPwd = value
	} else if key == strings.ToLower("-ta") {
		s.ta = value
	} else if key == strings.ToLower("-template") {
		s.template = value
	} else if key == strings.ToLower("-out") {
		s.OutputFolder = value
	}
}

func (s *Args) ShowHelp(folderPath string) {
	s.showLine("  -help:", "[可选]显示帮助")
	s.showLine("  -cfg:", fmt.Sprintf("[可选]指定配置文件路径, 默认: %s", filepath.Join(folderPath, "cfg", "ovpn.json")))
	s.showLine("  -out:", fmt.Sprintf("[可选]指定输出文件夹路径, 默认: %s", filepath.Join(folderPath, "clients")))
	s.showLine("  -ca.crt:", "[可选]指定CA证书文件路径")
	s.showLine("  -ca.key:", "[可选]指定CA密钥文件路径")
	s.showLine("  -ca.key.pwd:", "[可选]指定CA密钥文件密码")
	s.showLine("  -ta:", "[可选]指定TA(tls-auth)文件路径")
	s.showLine("  -template:", "[可选]指定客户端配置模板文件路径")
}

func (s *Args) showLine(label, value string) {
	fmt.Printf("%-15s %s", label, value)
	fmt.Println("")
}
