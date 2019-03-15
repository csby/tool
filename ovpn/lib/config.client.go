package lib

type ConfigClient struct {
	// 输出目录
	Folder string `json:"folder"`
	// 用户名
	OU string `json:"ou"`
	// 姓名
	CN string `json:"cn"`
	// 地区
	Locality string `json:"locality"`
	// 省份
	Province string `json:"province"`
	// 地址
	StreetAddress string `json:"address"`
	// 有效期(默认365)
	ExpiredDays int64 `json:"days"`
}
