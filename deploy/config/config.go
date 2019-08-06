package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	mutex sync.RWMutex

	Version     string `json:"version" note:"版本号"`
	Source      bool   `json:"source" note:"是否打包源代码"`
	Destination string `json:"destination" note:"输出目录"`

	Apps []App `json:"apps"`
}

func (s *Config) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) Example(folderPath string) {
	s.Version = "1.0.1.0"
	s.Source = false
	s.Destination = filepath.Join(folderPath, "rel")

	s.Apps = []App{
		{
			Enable: false,
			Name:   "vsgw",
			Bin: Binary{
				Root:  filepath.Join(folderPath, "bin"),
				Files: s.binaryFilesForGateway(),
			},
			Src: Source{
				Enable: true,
				Root:   filepath.Join(filepath.Dir(folderPath), "src", "github.com", "csby", "vsgw"),
				Ignore: []string{
					"tool",
					".git",
					".idea",
					".gitignore",
					"go.sum",
				},
			},
			Webs: []Web{
				{
					Enable: false,
					Name:   "doc",
					Src: Source{
						Root: filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(folderPath))), "vue", "doc"),
						Ignore: []string{
							"node_modules",
							"dist",
							".git",
							".idea",
							".gitignore",
							"README.md",
						},
					},
				},
				{
					Enable: false,
					Name:   "opt",
					Src: Source{
						Root: filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(folderPath))), "vue", "vsgw", "gateway", "opt"),
						Ignore: []string{
							"node_modules",
							"dist",
							".git",
							".idea",
							".gitignore",
							"README.md",
						},
					},
				},
			},
		},
		{
			Enable: s.enableAppForCrtMgr(),
			Name:   "certmaker",
			Bin: Binary{
				Root:  filepath.Join(folderPath, "bin"),
				Files: s.binaryFilesForCrtMgr(),
			},
			Src: Source{
				Enable: false,
			},
		},
		{
			Enable: true,
			Name:   "sqldm",
			Bin: Binary{
				Root:  filepath.Join(folderPath, "bin"),
				Files: s.binaryFilesForSlqDM(),
			},
			Src: Source{
				Enable: false,
			},
		},
	}
}
