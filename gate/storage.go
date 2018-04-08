package gate

import (
	"fmt"
	"github.com/dming/lodos/log"
	"github.com/dming/lodos/gate"
)

/**
存储用户的Session信息.
Session Bind Userid以后每次设置 settings都会调用一次Storage
*/
func (gate *Gate) Storage(Userid string, session gate.Session) (err error) {
	/*
	//增量更新settings
	var result *safemap.BeeMap4String
	if gate.storage.Check(Userid) {
		result = gate.storage.Get(Userid)
	} else {
		result = safemap.NewBeeMap4String()
	}
	for k, v := range settings {
		if !result.Check(k) {
			result.Set(k, v)
		}
	}*/
	//gate.storage.Set(Userid, session)

	log.Info("处理对Session的持久化完毕")
	return nil
}

/**
强制删除Session信息
*/
func (gate *Gate) Delete(Userid string) (err error) {
	log.Debug("删除Session持久化数据")
	//gate.storage.Delete(Userid)
	return nil
}

/**
获取用户Session信息
用户登录以后会调用Query获取最新信息
*/
func (gate *Gate) Query(Userid string) (data []byte, err error) {
	log.Debug("查询Session持久化数据")
/*
	if gate.storage.Check(Userid){
		return gate.storage.Get(Userid).Items(), nil
	}*/
	return nil, fmt.Errorf("can not find the storage [%s] settings ", Userid)
}

/**
用户心跳,一般用户在线时60s发送一次
可以用来延长Session信息过期时间
*/
func (gate *Gate) Heartbeat(Userid string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Gate Storage HeartBeat Error : %s", r)
		}
	}()
	log.Info("用户在线的心跳包")
/*
	if gate.storage.Check(Userid) {
		log.Debug("Gate storage :: Heartbeat")
		var tempSettings *safemap.BeeMap4String
		tempSettings = gate.storage.Get(Userid)

		if timeout := tempSettings.Get("timeout"); timeout != "" {
			if i, err := strconv.Atoi(timeout); err == nil && i < 120 {
				log.Info("更新用户 %s 在线的心跳包, now is %d", Userid, i + 60)
				tempSettings.Set("timeout", string(i + 60))
			}
		}
		gate.storage.Set(Userid, tempSettings)
	} else {
		log.Warning("can not find the storage user :  [%s] ", Userid)
	}
*/
}
