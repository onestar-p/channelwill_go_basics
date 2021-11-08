/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 17:51:37
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-05 15:35:33
 */
package config

import "time"

type ApplicationConfig struct {
	Name         string         `mapstructure:"name"`
	Network      string         `mapstructure:"network"`
	Ip           string         `mapstructure:"Ip"`
	Port         int            `mapstructure:"port"`
	HttpPort     int            `mapstructure:"httpport"`
	Env          string         `mapstructure:"env"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul"`
	DatabaseInfo DatabaseConfig `mapstructure:"database"`
	JwtInfo      JwtConfig      `mapstructure:"jwt"`
}

type ConsulConfig struct {
	Ip   string   `mapstructure:"Ip"`
	Port int      `mapstructure:"port"`
	Tags []string `mapstructure:"tags"`
}
type DatabaseConfig struct {
	Address  string `mapstructure:"address"`
	Config   string `mapstructure:"config"`
	DbType   string `mapstructure:"dbtype"`
	DbName   string `mapstructure:"dbname"`
	UserName string `mapstructure:"username"`
	PassWord string `mapstructure:"password"`
}

type JwtConfig struct {
	Issuer string        `mapstructure:"issuer"`
	Expire time.Duration `mapstructure:"expire"`
}
