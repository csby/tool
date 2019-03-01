package config

type App struct {
	Enable bool   `json:"enable"`
	Name   string `json:"name"`
	Bin    Binary `json:"bin"`
	Src    Source `json:"src"`
	Webs   []Web  `json:"webs"`
}
