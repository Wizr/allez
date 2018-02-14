package yezi_v4

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vettu/allez/libs/errorf"
	"github.com/vettu/allez/libs/middleware"
)

const (
	redisKeyCookie = "yezi-cookie"
	redisKeyPass   = "yezi-pass"
	durationCookie = time.Hour
	durationPass   = 5 * time.Minute
)

var id *Identity

func init() {
	id = &Identity{
		serialNumber: "C02LSLIPFH00",
		seller:       "yezi",
	}
}

func getRedisConn(c *gin.Context) (redisConn *redis.Client, err error) {
	if c, exist := c.Get(middleware.CtxRedisConn); !exist {
		err = errorf.New("no redis middleware or init failed")
	} else if c, ok := c.(*redis.Client); !ok {
		err = errorf.New("redis client type error")
	} else {
		redisConn = c
	}
	return
}

func GetAccounts(c *gin.Context) {
	var accounts []*sAccountInfo
	defer func() {
		if accounts == nil {
			c.JSON(http.StatusOK, gin.H{"error": "No result found."})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": accounts})
		}
	}()

	redisConn, err := getRedisConn(c)
	if err != nil {
		log.Printf("[yezi] getRedisConn | %v\n", err)
		return
	}
	doGetRedisValue(&passRedisKey{
		baseRedisKey: baseRedisKey{
			Conn:     redisConn,
			Key:      redisKeyPass,
			Duration: durationPass,
		},
	}, func(redisKey iRedisKey) {
		accounts = redisKey.(*passRedisKey).accounts
	}, func(redisKey iRedisKey) {
		accounts, err = doGetPasses(redisConn)
		if err != nil {
			log.Printf("[yezi] doGetPasses | %v\n", err)
			return
		}
		redisKey.(*passRedisKey).accounts = accounts
	})
}

func doGetRedisValue(redisKey iRedisKey, onSucc func(iRedisKey), onFail func(iRedisKey)) {
	var err error
	if err = redisKey.Get(); err != nil {
		log.Printf("[yezi] redis get %v: %v\n", redisKey.GetKey(), err)
	} else if err = redisKey.Deserialize(); err != nil {
		log.Printf("[yezi] redis deserialize %v: %v\n", redisKey.GetKey(), err)
	}
	if err != nil {
		onFail(redisKey)
		redisKey.Serialize()
		redisKey.Set()
		return
	}
	onSucc(redisKey)
}

func doGetPasses(redisConn *redis.Client) (accounts []*sAccountInfo, err error) {
	var cookies []*http.Cookie
	doGetRedisValue(&cookieRedisKey{
		baseRedisKey: baseRedisKey{
			Conn:     redisConn,
			Key:      redisKeyCookie,
			Duration: durationCookie,
		},
	}, func(redisKey iRedisKey) {
		cookies = redisKey.(*cookieRedisKey).cookies
	}, func(redisKey iRedisKey) {
		cookies, err = id.fetchToken()
		redisKey.(*cookieRedisKey).cookies = cookies
	})
	if err != nil {
		return
	}

	data, err := id.fetchAccountData(cookies)
	if err != nil {
		return
	}
	accounts, err = id.parseData(data)
	return
}
