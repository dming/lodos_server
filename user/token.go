package user

import (
	"github.com/dming/lodos/utils"
	"time"
)

const TOKEN_EXPIRED_SECOND int64 = int64(time.Hour) * 3 /1000 / 1000 /1000

var tokenCacheMap *TokenStatusFactory
func GetTokenCacheMap() *TokenStatusFactory {
	if tokenCacheMap == nil {
		tokenCacheMap = NewTokenMap(time.Second * 30)
	}
	return tokenCacheMap
}

func NewTokenMap(expired time.Duration) *TokenStatusFactory {
	t := &TokenStatusFactory{
		tokens: utils.NewBeeMap(),
	}
	t.OnInit()
	return t
}

type TokenStatusFactory struct {
	checkTimeout time.Duration
	tokens *utils.BeeMap //map[string]*TokenStatus
}

func (m *TokenStatusFactory) OnInit() {
	if m.checkTimeout > 0 {
		go m.timeout_handler()
	}
}

func (m *TokenStatusFactory) GetTokenStatus(token string) *TokenStatus {
	if t := m.tokens.Get(token); t != nil {
		if re, ok := t.(*TokenStatus); ok {
			return re
		}
	}
	return nil
}

func (m *TokenStatusFactory) SetTokenStatus(token string, status *TokenStatus) bool {
	return m.tokens.Set(token, status)
}

func (m *TokenStatusFactory) timeout_handler() {
	//todo: check timeout
}

//todo:
//1. get set update(from redis) token [get set finish]
//2. timeout_handler function


type TokenStatus struct {
	Token string
	Username string
	Login bool
	Expired int64
}