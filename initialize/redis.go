package initialize

import (
	"fmt"

	"github.com/go-redis/redis"

	"channelwill_go_basics/dao"
	"channelwill_go_basics/global"
)

func InitRedis() {
	redisConfig := global.ApplicationConfig.RedisInfo
	dao.Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s", redisConfig.Address),
	})
}
