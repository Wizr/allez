package yezi_v4

import (
	"net/http"
	"strings"

	"github.com/vettu/allez/libs/errorf"
)

type cookieRedisKey struct {
	baseRedisKey
	cookies []*http.Cookie
}

func (crk *cookieRedisKey) Serialize() {
	var accounts []string
	for _, ai := range crk.cookies {
		accounts = append(accounts, ai.Name+","+ai.Value)
	}
	crk.rawValue = strings.Join(accounts, ";")
}

func (crk *cookieRedisKey) Deserialize() (err error) {
	cookies := strings.Split(crk.rawValue, ";")
	for _, c := range cookies {
		kv := strings.Split(c, ",")
		if len(kv) != 2 {
			continue
		}
		crk.cookies = append(crk.cookies, &http.Cookie{
			Name:  kv[0],
			Value: kv[1],
		})
	}
	if cookies == nil {
		err = errorf.Newf("Deserialize failed for key %v", crk.Key)
	}
	return
}
