package main

import (
	"bufio"
	"fmt"
	"github.com/csby/tool/module"
	"github.com/csby/tool/ovpn/lib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	moduleType    = "tool"
	moduleName    = "ovpn-cmd"
	moduleRemark  = "OpenVPN客户端证书工具"
	moduleVersion = "1.0.1.0"
)

var (
	args = &Args{}
	cfg  = &lib.Config{}
)

func init() {
	moduleArgs := &module.ArgsModule{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, args)
	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())

	cfgPath := args.cfg
	if cfgPath == "" {
		cfgPath = filepath.Join(rootFolder, "cfg", "ovpn.json")
	} else if !filepath.IsAbs(cfgPath) {
		cfgPath = filepath.Join(rootFolder, cfgPath)
	}
	fmt.Println("cfg:", cfgPath)

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		cfg.Example()
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

	if args.help {
		args.ShowHelp(rootFolder)
		os.Exit(0)
	}

	if len(args.caCrt) > 0 {
		v, err := readFile(args.caCrt)
		if err != nil {
			fmt.Println("read ca cert file fail:", err)
		} else {
			cfg.Ca.Cert = v
			cfg.SaveToFile(cfgPath)
		}
	}
	if len(args.caKey) > 0 {
		v, err := readFile(args.caKey)
		if err != nil {
			fmt.Println("read ca key file fail:", err)
		} else {
			cfg.Ca.Key = v
			cfg.SaveToFile(cfgPath)
		}
	}
	if len(args.caKeyPwd) > 0 {
		if cfg.Ca.KeyPassword != args.caKeyPwd {
			cfg.Ca.KeyPassword = args.caKeyPwd
			cfg.SaveToFile(cfgPath)
		}
	}
	if len(args.ta) > 0 {
		v, err := readFile(args.ta)
		if err != nil {
			fmt.Println("read ta file fail:", err)
		} else {
			cfg.Ta = v
			cfg.SaveToFile(cfgPath)
		}
	}
	if len(args.template) > 0 {
		v, err := readFile(args.template)
		if err != nil {
			fmt.Println("read template file fail:", err)
		} else {
			cfg.Template = v
			cfg.SaveToFile(cfgPath)
		}
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入用户名[%s]:", cfg.Client.OU)
	input, err := readInput(inputReader)
	if err == nil {
		if cfg.Client.OU != input {
			cfg.Client.OU = input
			cfg.SaveToFile(cfgPath)
		}
	}

	if len(cfg.Client.CN) < 1 {
		cfg.Client.CN = cfg.Client.OU
	}
	fmt.Printf("请输入姓名[%s]:", cfg.Client.CN)
	input, err = readInput(inputReader)
	if err == nil {
		if cfg.Client.CN != input {
			cfg.Client.CN = input
			cfg.SaveToFile(cfgPath)
		}
	}

	fmt.Printf("请输入地区[%s]:", cfg.Client.Locality)
	input, err = readInput(inputReader)
	if err == nil {
		if cfg.Client.Locality != input {
			cfg.Client.Locality = input
			cfg.SaveToFile(cfgPath)
		}
	}

	fmt.Printf("请输入省份[%s]:", cfg.Client.Province)
	input, err = readInput(inputReader)
	if err == nil {
		if cfg.Client.Province != input {
			cfg.Client.Province = input
			cfg.SaveToFile(cfgPath)
		}
	}

	fmt.Printf("请输入地址[%s]:", cfg.Client.StreetAddress)
	input, err = readInput(inputReader)
	if err == nil {
		if cfg.Client.StreetAddress != input {
			cfg.Client.StreetAddress = input
			cfg.SaveToFile(cfgPath)
		}
	}

	if cfg.Client.ExpiredDays < 1 {
		cfg.Client.ExpiredDays = 365
	}
	fmt.Printf("有效期(天)[%d]:", cfg.Client.ExpiredDays)
	input, err = readInput(inputReader)
	if err == nil {
		inputVal, err := strconv.ParseInt(input, 10, 64)
		if err == nil && inputVal > 0 {
			if cfg.Client.ExpiredDays != inputVal {
				cfg.Client.ExpiredDays = inputVal
				cfg.SaveToFile(cfgPath)
			}
		}
	}

	if len(args.OutputFolder) > 0 {
		cfg.Client.Folder = args.OutputFolder
	}
	if len(cfg.Client.Folder) < 1 {
		cfg.Client.Folder = filepath.Join(rootFolder, "clients")
	}
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	inputVal := strings.Replace(input, "\n", "", -1)
	if len(inputVal) > 0 {
		return inputVal, nil
	}

	return "", fmt.Errorf("empty")
}
