package main

import (
	"fmt"
	"github.com/csby/tool/ovpn/lib"
	"os"
)

func main() {
	if len(cfg.Client.OU) < 1 {
		fmt.Println("错误: 用户名为空")
		os.Exit(0)
	}

	ovpn := &lib.OVpn{}
	path, err := ovpn.CreateClientConf(cfg)
	if err != nil {
		fmt.Println("失败:", err)
	} else {
		fmt.Println("成功:", path)
	}
}
