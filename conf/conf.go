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
	viper.Reset()
	// 设置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "a12bCd3_W45pUq6")
	viper.SetDefault("database.dbname", "mysql")

	viper.SetConfigName("ghl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../../conf")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认值
		fmt.Printf("配置文件未找到，使用默认配置: %v\n", err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("配置反序列化失败: %w", err))
	}
	return c
}
