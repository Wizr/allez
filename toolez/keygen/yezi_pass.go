package keygen

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vettu/allez/libs"
)

type AccountRequestInfo struct {
	ApplicationID string
	ClientID      string
	AES128Key     string
	URL           string
}
type AccountInfo struct {
	AppleID  string
	Password string
}

var yeziAccounts = map[string]*AccountRequestInfo{
	"yezi": {
		ApplicationID: "0jPhBkJiLAPJNbaYFBAJavJR-gzGzoHsz",
		ClientID:      "T2EcwHL2EjWR46PNNGytKEI8",
		AES128Key:     "i know nothing about yezi!KALSDFIOQPWREJ91203JVZLVKJ0-1234",
		URL:           "https://api.leancloud.cn/1.1/classes/Account/570b9b505bbb50004c1c6aee",
	},
	"yezi3": {
		ApplicationID: "mywR1aPpOqPIOnrtJSSPWVzk-gzGzoHsz",
		ClientID:      "pxiwVFBgLhjC1pHMFLSB8sqr",
		AES128Key:     "i know nothing about yezi!KALSDFIOQPWREJ91203JVZLVKJ0-1234",
		URL:           "https://api.leancloud.cn/1.1/classes/Account/56ceaa83c24aa80052109e07",
	},
}

func GetAccounts(c *gin.Context) {
	result := map[string]*AccountInfo{}
	for name, info := range yeziAccounts {
		if accountInfo := GetAccountInfo(info); accountInfo != nil {
			result[name] = accountInfo
		}
	}
	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No result found."})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func GetAccountInfo(requestInfo *AccountRequestInfo) *AccountInfo {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("GET", requestInfo.URL, nil)
	req.Header.Add("User-Agent", "AVOS Cloud OSX-v3.2.2 SDK")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-LC-Id", requestInfo.ApplicationID)
	req.Header.Add("X-LC-Sign", X_LS_Sign(requestInfo.ClientID))
	req.Header.Add("X-LC-Prod", "1")
	// do request
	if resp, err := client.Do(req); err != nil {
		log.Println(err)
		return nil
	} else {
		defer resp.Body.Close()
		// read response body
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
			return nil
		} else {
			// parse response info
			accountInfo := &AccountInfo{}
			if json.Unmarshal(body, accountInfo) != nil {
				log.Println(err)
				return nil
			}
			key := []byte(requestInfo.AES128Key)[:32]
			if block, err := aes.NewCipher(key); err != nil {
				log.Println(err)
				return nil
			} else {
				// decode base64 encoded password
				passBase64 := accountInfo.Password
				if pass, err := base64.StdEncoding.DecodeString(passBase64); err != nil {
					log.Println(err)
					return nil
				} else {
					// decrypt aes ecb encoded password
					// https://github.com/dev5tec/FBEncryptor
					bm := libs.NewECBDecrypter(block)
					passOut := make([]byte, len(pass))
					bm.CryptBlocks(passOut, pass)
					accountInfo.Password = string(passOut)
					return &AccountInfo{
						AppleID:  accountInfo.AppleID,
						Password: accountInfo.Password,
					}
				}
			}
		}
	}
}

func X_LS_Sign(ClientID string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	md5In := timestamp + ClientID
	md5Result := md5.Sum([]byte(md5In))
	md5Out := hex.EncodeToString(md5Result[:])
	sig := fmt.Sprintf("%s,%s", md5Out, timestamp)
	return sig
}
