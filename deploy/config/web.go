package config

type Web struct {
	Enable bool   `json:"enable"`
	Name   string `json:"name"`
	Src    Source `json:"src"`
}
