/*
@Time : 2020/8/20 21:07
@Author : liangjiefan
*/
package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	onceDo sync.Once
)

type Config struct {
	Server *ServerConfig
	Mysql  *MysqlConfig
}

type ServerConfig struct {
	Port    int
	RunMode string
}

type MysqlConfig struct {
	Addr        string
	Db          string
	User        string
	Password    string
	MaxIdleCons int
	MaxOpenCons int
}

func ProvideConfig(path string) (*Config, error) {
	c := new(Config)
	onceDo.Do(func() {
		c.initConfig(path)
		c.init()
	})

	return c, nil
}

func (c *Config) initConfig(path string) {
	if path == "" {
		path = "."
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	viper.SetConfigType("yaml")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return
}

func (c *Config) init() {
	c.Mysql = &MysqlConfig{
		Addr:        viper.GetString("mysql.addr"),
		Db:          viper.GetString("mysql.db"),
		User:        viper.GetString("mysql.user"),
		Password:    viper.GetString("mysql.password"),
		MaxIdleCons: viper.GetInt("mysql.maxIdleCons"),
		MaxOpenCons: viper.GetInt("mysql.maxOpenCons"),
	}

	c.Server = &ServerConfig{
		Port:    viper.GetInt("server.port"),
		RunMode: viper.GetString("server.runmode"),
	}
}
