package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Ca  ConfigCa  `json:"ca"`
	Crt ConfigCrt `json:"crt"`
	Crl ConfigCrl `json:"crl"`
}

func NewConfig() *Config {
	return &Config{
		Crt: ConfigCrt{
			Organization: "client",
			ExpiredDays:  365,
			Hosts:        make([]string, 0),
		},
	}
}

func (s *Config) SaveToFile(filePath string) error {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

func (s *Config) LoadFromFile(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}
