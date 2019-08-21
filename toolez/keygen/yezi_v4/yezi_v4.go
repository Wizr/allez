package yezi_v4

import (
	"log"
	"net/http"
	"time"

	"github.com/Wizr/allez/libs/errorf"
	"github.com/Wizr/allez/libs/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	}, func(redisKey iRedisKey) error {
		accounts, err = doGetPasses(redisConn)
		if err != nil {
			log.Printf("[yezi] doGetPasses | %v\n", err)
			redisConn.Del(redisKeyCookie)
			redisConn.Del(redisKeyPass)
			return err
		}
		redisKey.(*passRedisKey).accounts = accounts
		return nil
	})
}

func doGetRedisValue(redisKey iRedisKey, onSucc func(iRedisKey), onFail func(iRedisKey) error) {
	var err error
	if err = redisKey.Get(); err != nil {
		log.Printf("[yezi] redis get %v: %v\n", redisKey.GetKey(), err)
	} else if err = redisKey.Deserialize(); err != nil {
		log.Printf("[yezi] redis deserialize %v: %v\n", redisKey.GetKey(), err)
	}
	if err != nil {
		if onFail(redisKey) == nil {
			redisKey.Serialize()
			redisKey.Set()
		}
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
	}, func(redisKey iRedisKey) error {
		cookies, err = id.fetchToken()
		if err != nil {
			return err
		}
		redisKey.(*cookieRedisKey).cookies = cookies
		return nil
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
