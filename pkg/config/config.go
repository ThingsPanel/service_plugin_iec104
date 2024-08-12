package config

type Config struct {
	Things Things `yaml:"things"`
	Logs   Logs   `yaml:"logs"`
}

type Things struct {
	Local  string `yaml:"local"`  // 本地服务监听地址
	Remote string `yaml:"remote"` // thingspanel api地址
	Mqtt   Mqtt   `yaml:"mqtt"`   // thingspanel mqtt配置
	Device Device `yaml:"device"`
}

type Mqtt struct {
	Broker   string `yaml:"broker"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Device struct {
	Timeout int `yaml:"timeout"` // 设备数据超时时间, 超时后判为离线, 单位（分钟）
}

type Logs struct {
	Debug bool   `yaml:"debug"` // 是否开启debug
	File  string `yaml:"file"`  // 日志文件保存位置
}
