package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Args struct {
	cfg  string
	ver  string
	out  string
	src  *bool
	help bool
}

func (s *Args) Parse(key, value string) {
	if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.help = true
	} else if key == strings.ToLower("-cfg") {
		s.cfg = value
	} else if key == strings.ToLower("-ver") {
		s.ver = value
	} else if key == strings.ToLower("-out") {
		s.out = value
	} else if key == strings.ToLower("-src") {
		src := true
		if strings.ToLower(value) == "false" {
			src = false
		}
		s.src = &src
	}
}

func (s *Args) ShowHelp(folderPath string) {
	s.showLine("  -help:", "[可选]显示帮助")
	s.showLine("  -cfg:", fmt.Sprintf("[可选]指定配置文件路径, 默认: %s", filepath.Join(folderPath, "cfg", "deploy.json")))
	s.showLine("  -ver:", "[可选]指定版本号(如果指定将覆盖配置文件中版本号), 格式:major.minor.build.revision, 如-ver=1.0.8.2, 默认: 1.0.1.0")
	s.showLine("  -out:", fmt.Sprintf("[可选]指定输出根目录, 默认: %s", filepath.Join(folderPath, "rel")))
	s.showLine("  -src:", "[可选]打包源代码(如果指定将覆盖配置文件中的source项), 如-src=true")
}

func (s *Args) showLine(label, value string) {
	fmt.Printf("%-8s %s", label, value)
	fmt.Println("")
}
