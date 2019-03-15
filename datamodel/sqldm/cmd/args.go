package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Args struct {
	cfg  string
	help bool
}

func (s *Args) Parse(key, value string) {
	if key == strings.ToLower("-h") ||
		key == strings.ToLower("-help") ||
		key == strings.ToLower("--help") {
		s.help = true
	} else if key == strings.ToLower("-cfg") {
		s.cfg = value
	}
}

func (s *Args) ShowHelp(folderPath string) {
	s.showLine("  -help:", "[可选]显示帮助")
	s.showLine("  -cfg:", fmt.Sprintf("[可选]指定配置文件路径, 默认: %s", filepath.Join(folderPath, "sqldm.json")))
}

func (s *Args) showLine(label, value string) {
	fmt.Printf("%-8s %s", label, value)
	fmt.Println("")
}
