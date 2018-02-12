package keygen

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vettu/allez/libs"
	"github.com/vettu/allez/libs/errorf"
	"github.com/vettu/allez/libs/middleware"
)

const (
	tokenUrl       = "https://api.yeziapp.com/client/tokens"
	accountUrl     = "https://api.yeziapp.com/client/accounts"
	redisCookieKey = "yezi-cookie"
)

type accountInfo struct {
	AppleID  string
	Password string
}

var id *identity

func init() {
	id = &identity{
		serialNumber: "C02LSLIPFH00",
		seller:       "yezi",
	}
}

func GetAccounts(c *gin.Context) {
	accounts := getAccountInfo(c)
	if accounts == nil {
		c.JSON(http.StatusOK, gin.H{"error": "No result found."})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": accounts})
	}
}

func getAccountInfo(c *gin.Context) (accounts []*accountInfo) {
	var redisConn *redis.Client
	if c, exist := c.Get(middleware.CtxRedisConn); !exist {
		log.Println("[yezi] no redis middleware or init failed")
	} else if c, ok := c.(*redis.Client); !ok {
		log.Println("[yezi] redis client type error")
	} else {
		redisConn = c
	}
	if redisConn == nil {
		c.JSON(http.StatusOK, gin.H{"error": "No result found."})
		return
	}
	// get cookie
	var cookie []*http.Cookie
	rawCookie, err := redisConn.Get(redisCookieKey).Result()
	if err == redis.Nil {
		// fetch cookie
		cookie, err = id.fetchToken()
		if err != nil {
			log.Printf("[yezi] fetchToken | %v\n", err)
			return
		}
		// cache cookie in redis
		var rawCookies []string
		var expire *time.Time
		for _, c := range cookie {
			rawCookies = append(rawCookies, fmt.Sprintf("%v,%v", c.Name, c.Value))
			if expire == nil {
				expire = &c.Expires
			} else if !c.Expires.IsZero() && c.Expires.Before(*expire) {
				expire = &c.Expires
			}
		}
		var duration time.Duration
		if expire != nil {
			duration = expire.Sub(time.Now()) - time.Hour
		}
		rawCookie = strings.Join(rawCookies, ";")
		err = redisConn.Set(redisCookieKey, rawCookie, duration).Err()
		if err != nil {
			log.Printf("[yezi] redis set error | %v\n", err)
		}
	} else if err != nil {
		log.Printf("[yezi] redis get error | %v\n", err)
		return
	}
	// get account data
	cookies := strings.Split(rawCookie, ";")
	for _, c := range cookies {
		kv := strings.Split(c, ",")
		if len(kv) != 2 {
			continue
		}
		cookie = append(cookie, &http.Cookie{
			Name:  kv[0],
			Value: kv[1],
		})
	}
	data, err := id.fetchAccountData(cookie)
	if err != nil {
		if err := redisConn.Del(redisCookieKey).Err(); err != nil {
			log.Printf("[yezi] redis del error | %v\n", err)
		}
		log.Printf("[yezi] fetchAccountData | %v\n", err)
		return
	}
	accounts, err = id.parseData(data)
	if err != nil {
		log.Printf("[yezi] parseData | %v\n", err)
		return nil
	}
	return
}

type identity struct {
	serialNumber string
	seller       string
}

func (id *identity) fetchToken() (cookie []*http.Cookie, erf error) {
	jsonStr := []byte (fmt.Sprintf(`{"serialNumber": "%v", "seller": "%v"}`, id.serialNumber, id.seller))
	resp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		erf = errorf.Newf("Post error: %v", err)
		return
	}
	defer resp.Body.Close()
	cookie = resp.Cookies()
	if len(cookie) == 0 {
		erf = errorf.Newf("No cookie found: %v", err)
		return
	}
	return
}

func (id *identity) fetchAccountData(cookie []*http.Cookie) (data string, erf error) {
	// create request
	req, _ := http.NewRequest("GET", accountUrl, nil)
	req.Close = true
	for _, c := range cookie {
		req.AddCookie(c)
	}

	// do request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		erf = errorf.Newf("do request failed: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		erf = errorf.Newf("request not succeeds", resp.StatusCode)
		return
	}

	// read response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		erf = errorf.Newf("read response body failed: %v", err)
		return
	}

	// parse response
	t := &struct{ Data string }{}
	err = json.Unmarshal(body, t)
	if err != nil {
		erf = errorf.Newf("json unmarshal failed: %v", err)
		return
	}
	data = t.Data
	return
}

// decrypt crypto-js encrypted data
func (id *identity) parseData(data string) (accounts []*accountInfo, erf error) {
	// prepare encrypted data
	rawData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		erf = errorf.Newf("base64 decoding failed %v", err)
		return
	}

	// generate and derive key
	rs := fmt.Sprintf("%v@%v", id.serialNumber, id.seller)
	pass := md5.Sum([]byte(rs))
	passHex := make([]byte, hex.EncodedLen(len(pass)))
	hex.Encode(passHex, pass[:])
	// will be 48 bytes, first 32 bytes for creating block, last 16 is IV
	key := libs.EVPKDF(passHex, rawData[8:16], 1, 12, func() hash.Hash {
		return md5.New()
	})

	// decrypt data with the derived key
	block, _ := aes.NewCipher(key[:32])
	cbc := cipher.NewCBCDecrypter(block, key[32:])
	d := rawData[block.BlockSize():]
	cbc.CryptBlocks(d, d)
	d = libs.PKCS5UnPadding(d)

	// parse the decrypted data, a json
	err = json.Unmarshal(d, &accounts)
	if err != nil {
		erf = errorf.Newf("json unmarshal failed: %v", string(d))
		return
	}
	return
}
