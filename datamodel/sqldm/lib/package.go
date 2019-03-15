package lib

type Package struct {
	Name   string `json:"name" note:"名称，如: data"`
	Path   string `json:"path" note:"路径，如: github.com/project/data"`
	Folder string `json:"folder" note:"文件夹，如: /home/user/go/src/github.com/project/data"`
}
