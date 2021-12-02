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
	Name          string         `mapstructure:"name"`
	Network       string         `mapstructure:"network"`
	Ip            string         `mapstructure:"Ip"`
	Port          int            `mapstructure:"port"`
	AutoPort      bool           `mapstructure:"autoport"`
	HttpPort      int            `mapstructure:"httpport"`
	Env           string         `mapstructure:"env"`
	MachineId     uint16         `mapstructure:"machine_id"`
	AuthSrv       string         `mapstructure:"auth_srv"`
	EtranslateSrv string         `mapstructure:"etranslate_srv"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul"`
	DatabaseInfo  DatabaseConfig `mapstructure:"database"`
	JwtInfo       JwtConfig      `mapstructure:"jwt"`
	AESInfo       AesConfig      `mapstructure:"aes"`
	EmailInfo     EmailConfig    `mapstructure:"email"`
	RedisInfo     RedisConfig    `mapstructure:"redis"`
	DingtalkInfo  DingtalkConfig `mapstructure:"dingtalk"`
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

type AesConfig struct {
	Key  string `mapstructure:"key"`
	Iv   string `mapstructure:"iv"`
	Mode string `mapstructure:"mode"`
}

type EmailConfig struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"username"`
	Fromname string `mapstructure:"fromname"`
	Passwd   string `mapstructure:"passwd"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"address"`
	DB       int    `mapstructure:"db"`
}

type DingtalkConfig struct {
	AccessToken string `mapstructure:"accesstoken"`
}
