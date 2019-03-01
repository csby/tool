package main

import (
	"fmt"
	"github.com/csby/tool/deploy/packer"
)

func main() {
	pack := packer.NewPacker(cfg)
	err := pack.Pack()
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("成功,输出目录:", pack.OutputFolder())
	}
}
