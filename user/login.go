package user

import (
	"github.com/dming/lodos/gate"
	"fmt"
	"github.com/dming/lodos/db/base"
	"github.com/globalsign/mgo/bson"
	"github.com/dming/lodos/gate/base"
	"github.com/dming/lodos/log"
	"encoding/json"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/dming/lodos/rpc/base"
)

// return token ?
func (m *user) login(session gate.Session, msg map[string]interface{}) (token string, err error) {
	//time.Sleep(time.Millisecond * 200)
	if msg["username"] == nil || msg["password"] == nil {
		err = fmt.Errorf("userName or passWord cannot be nil")
		return "", baserpc.NewError(baserpc.ClientParamNotAdapted, err)
	}
	var username string
	if u, ok := msg["username"].(string); !ok {
		err = fmt.Errorf("username cannot be format as [string]")
		return "", baserpc.NewError(baserpc.ClientFormatError, err)
	} else {
		username = u
	}
	var password string
	if p, ok := msg["password"].(string); !ok {
		err = fmt.Errorf("password cannot be format as [string]")
		return "", baserpc.NewError(baserpc.ClientFormatError, err)
	} else {
		password = p
	}

	//TODO: check if valid of username and password, if valid, set token, else token remain to ""
	mongoInfo := m.GetModuleSetting().Mongo
	coll, err := basedb.GetMongoFactories().GetCOLL(mongoInfo.Uri, mongoInfo.DB, USERCOLL)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	result := struct {
		A,B int
	}{}
	err = coll.Find(bson.M{"username": username, "password": password}).One(&result)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	} else {
		token, err = m.createToken(username)
		if err != nil {
			return "", baserpc.NewError(baserpc.ServerDBError, err)
		}
	}

	if token == "" {
		return "", fmt.Errorf("login fail with %s and %s", username, password)
	}
	err = session.Bind(username)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerRpcInvokeError, err)
	}
	err = session.Push() //推送到网关
	if err != nil {
		log.Error("cannot push , err is %s", err.Error())
		return "", baserpc.NewError(baserpc.ServerRpcInvokeError, err)
	}

	//session.SetToken(token)
	//todo: save the token to redis server db
	m.redisConn = basedb.GetRedisFactory().GetPool(m.GetModuleSetting().Redis.DBUri).Get()
	if m.redisConn.Err() != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, m.redisConn.Err())
	}
	//expired := (time.Now().UTC().UnixNano() + TOKEN_EXPIRED) / 1000/ 1000 / 1000 //=> second
	tValue := &TokenStatus{
		Token: token,
		Username: username,
		Login: true,
	}
	b, err := json.Marshal(tValue)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerFormatError, err)
	}

	re, err := redis.String(m.redisConn.Do("Get", username))//, {})
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	_, err = m.redisConn.Do("Del", re)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	_, err = m.redisConn.Do("Set", username, token)//, {})
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	_, err = m.redisConn.Do("Set", token, b)//, {})
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	_, err = m.redisConn.Do("Expire", token, TOKEN_EXPIRED_SECOND)
	if err != nil {
		return "", baserpc.NewError(baserpc.ServerDBError, err)
	}
	tValue.Expired = (time.Now().UTC().UnixNano() + TOKEN_EXPIRED_SECOND * int64(time.Second)) / int64(time.Millisecond)
	GetTokenCacheMap().SetTokenStatus(token, tValue)
	return token, nil
}

func (m *user) authentication(session gate.Session, username string, token string) (bool, error) {
	//TODO: authenticate the username and token in the session
	if session.GetToken() == token && session.GetUserid() == username {
		return true, nil
	}
	return false, nil
}

func (m *user) createToken(username string) (string, error) {
	//todo: create token by given username
	return basegate.Get_uuid(), nil
}
