/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-05 00:07:05
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-05 15:09:22
 */
package client

import (
	"context"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"channelwill_go_basics/ent"
	"channelwill_go_basics/ent/migrate"
	"channelwill_go_basics/global"
)

var (
	Db    *ent.Client
	Redis *redis.Client
)

func InitMysql() {
	// MySQL
	DbConfig := global.ApplicationConfig.DatabaseInfo
	username := DbConfig.UserName
	password := DbConfig.PassWord
	dbname := DbConfig.DbName
	address := DbConfig.Address
	config := DbConfig.Config
	dsn := username + ":" + password + "@tcp(" + address + ")/" + dbname + "?" + config
	zap.S().Infof("dsn: %v", dsn)

	var err error
	Db, err = ent.Open(DbConfig.DbType, dsn)
	if err != nil {
		zap.S().Fatalf("failed opening connection to mysql: %v", err)
	}

	if err := Db.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		zap.S().Panicf("Mysql create schema err: %v", err)
	}

}
