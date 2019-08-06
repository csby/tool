package main

import (
	"fmt"
	"github.com/csby/tool/deploy/config"
	"github.com/csby/tool/module"
	"os"
	"path/filepath"
)

const (
	moduleType    = "tool"
	moduleName    = "deploy"
	moduleRemark  = "系统发布打包工具"
	moduleVersion = "1.0.1.0"
)

var (
	args = &Args{}
	cfg  = &config.Config{}
)

func init() {
	moduleArgs := &module.ArgsModule{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, args)
	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())

	cfgPath := args.cfg
	if cfgPath == "" {
		cfgPath = filepath.Join(rootFolder, "cfg", "deploy.json")
	} else if !filepath.IsAbs(cfgPath) {
		cfgPath = filepath.Join(rootFolder, cfgPath)
	}
	fmt.Println("cfg:", cfgPath)

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		cfg.Example(rootFolder)
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(cfgPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}
	//fmt.Println(cfg.FormatString())

	updateCount := 0
	if args.ver != "" {
		if args.ver != cfg.Version {
			cfg.Version = args.ver
			updateCount++
		}
	}
	if args.out != "" {
		if args.out != cfg.Destination {
			cfg.Destination = args.out
			updateCount++
		}
	}
	if args.src != nil {
		if *args.src != cfg.Source {
			cfg.Source = *args.src
			updateCount++
		}
	}

	if args.help {
		args.ShowHelp(rootFolder)
		os.Exit(0)
	}
	if cfg.Version == "" {
		fmt.Println("错误: 必须指定版本号")
		args.ShowHelp(rootFolder)
		os.Exit(0)
	}

	if updateCount > 0 {
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("update configure file fail: ", err)
		}
	}
}
