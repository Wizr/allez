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

	"github.com/gin-gonic/gin"
	"github.com/vettu/allez/libs"
	"github.com/vettu/allez/libs/errorf"
)

const (
	tokenUrl   = "https://api.yeziapp.com/client/tokens"
	accountUrl = "https://api.yeziapp.com/client/accounts"
	maxTry     = 2
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
	id.token, _ = id.fetchToken()
}

func GetAccounts(c *gin.Context) {
	accounts := func() (accounts []*accountInfo) {
		var err error
		if id.token == nil {
			if id.curTry == maxTry {
				return
			}
			id.token, err = id.fetchToken()
			id.curTry ++
			if err != nil {
				log.Printf("[yezi] fetchToken | %v", err)
				return
			}
		}
		id.data, err = id.fetchAccountData()
		if err != nil {
			log.Printf("[yezi] parseData | %v", err)
			return
		}
		accounts, err = id.parseData()
		if err != nil {
			log.Printf("[yezi] parseData | %v", err)
			return nil
		}
		return
	}()
	if accounts == nil {
		c.JSON(http.StatusOK, gin.H{"error": "No result found."})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": accounts})
	}
}

type identity struct {
	serialNumber string
	seller       string

	token  []*http.Cookie
	data   string
	curTry int
}

func (id *identity) fetchToken() (token []*http.Cookie, erf error) {
	jsonStr := []byte (fmt.Sprintf(`{"serialNumber": "%v", "seller": "%v"}`, id.serialNumber, id.seller))
	resp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		erf = errorf.Newf("Post error: %v\n", err)
		return
	}
	token = resp.Cookies()
	return
}

func (id *identity) fetchAccountData() (data string, erf error) {
	// create request
	req, _ := http.NewRequest("GET", accountUrl, nil)
	for _, cookie := range id.token {
		req.AddCookie(cookie)
	}

	// do request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		erf = errorf.Newf("do request failed: %v\n", err)
		return
	}

	// read response
	defer func() {
		resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		erf = errorf.Newf("read response body failed: %v\n", err)
		return
	}

	// parse response
	t := &struct{ Data string }{}
	err = json.Unmarshal(body, t)
	if err != nil {
		erf = errorf.Newf("json unmarshal failed: %v\n", err)
		return
	}
	data = t.Data
	return
}

// decrypt crypto-js encrypted data
func (id *identity) parseData() (accounts []*accountInfo, erf error) {
	// prepare encrypted data
	rawData, err := base64.StdEncoding.DecodeString(id.data)
	if err != nil {
		erf = errorf.Newf("base64 decoding failed %v\n", err)
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
		erf = errorf.Newf("json unmarshal failed: %v\n", string(d))
		return
	}
	return
}
