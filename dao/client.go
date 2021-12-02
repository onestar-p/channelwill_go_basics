/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-05 00:07:05
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-05 15:09:22
 */
package dao

import (
	"channelwill_go_basics/ent"

	"github.com/go-redis/redis"
)

var (
	Db    *ent.Client
	Redis *redis.Client
)
