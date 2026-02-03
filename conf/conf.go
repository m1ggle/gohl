package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func LoadConf() Config {
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("配置反序列化失败: %w", err))
	}
	return c
}

func InitViper() {
	// 设置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "mysql")

	viper.SetEnvPrefix("GOHL")
	viper.AutomaticEnv()

	// 绑定环境变量
	viper.BindEnv("database.password", "GOHL_DB_PASS")
	viper.BindEnv("database.user", "GOHL_DB_USER")
	viper.BindEnv("database.host", "GOHL_DB_HOST")
	viper.BindEnv("database.port", "GOHL_DB_PORT")
	viper.BindEnv("database.dbname", "GOHL_DB_NAME")
}
