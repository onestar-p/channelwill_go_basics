/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 18:31:20
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-05 15:09:38
 */
package initialize

import (
	"channelwill_go_basics/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	etRoot "channelwill_go_basics"
)

func InitConfig() {

	configFileName := etRoot.Path("config.yaml")
	v := viper.New()
	// 设置文件路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ApplicationConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息：%v", global.ApplicationConfig)

	// 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件产生变化：%s", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ApplicationConfig)
		zap.S().Infof("配置信息：%v", global.ApplicationConfig)
	})

}
