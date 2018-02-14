package yezi_v4

import (
	"strings"

	"github.com/vettu/allez/libs/errorf"
)

type sAccountInfo struct {
	AppleID  string
	Password string
}

type passRedisKey struct {
	baseRedisKey
	accounts []*sAccountInfo
}

func (prk *passRedisKey) Serialize() {
	var accounts []string
	for _, ai := range prk.accounts {
		accounts = append(accounts, ai.AppleID+","+ai.Password)
	}
	prk.rawValue = strings.Join(accounts, ";")
}

func (prk *passRedisKey) Deserialize() (err error) {
	accounts := strings.Split(prk.rawValue, ";")
	for _, c := range accounts {
		kv := strings.Split(c, ",")
		if len(kv) != 2 {
			continue
		}
		prk.accounts = append(prk.accounts, &sAccountInfo{
			AppleID:  kv[0],
			Password: kv[1],
		})
	}
	if accounts == nil {
		err = errorf.Newf("Deserialize failed for key %v", prk.Key)
	}
	return
}
