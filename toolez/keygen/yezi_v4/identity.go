package yezi_v4

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
	"net/http"

	"github.com/vettu/allez/libs"
	"github.com/vettu/allez/libs/errorf"
)

const (
	tokenUrl   = "https://api.yeziapp.com/client/tokens"
	accountUrl = "https://api.yeziapp.com/client/accounts"
)

type Identity struct {
	serialNumber string
	seller       string
}

func (id *Identity) fetchToken() (cookie []*http.Cookie, erf error) {
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

func (id *Identity) fetchAccountData(cookie []*http.Cookie) (data string, erf error) {
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
func (id *Identity) parseData(data string) (accounts []*sAccountInfo, erf error) {
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
