/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package gate

import (
	"fmt"
	"github.com/dming/lodos/conf"
	"github.com/dming/lodos/gate"
	log "github.com/dming/lodos/mlog"
	"github.com/dming/lodos/module"
	"github.com/DeanThompson/syncmap"
	"strconv"
)

var Module = func() module.Module {
	gate := new(Gate)
	return gate
}

type Gate struct {
	gate.Gate //继承

	//storage map[string]map[string]string
	storage *MutexMap
	storage2 *syncmap.SyncMap
}

func (gate *Gate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Gate"
}
func (gate *Gate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (gate *Gate) OnInit(app module.AppInterface, settings conf.ModuleSettings) {
	//注意这里一定要用 gate.Gate 而不是 module.BaseModule
	gate.Gate.OnInit(app, settings)

	gate.storage = NewMutexMap()
	gate.storage2 = syncmap.New()
	gate.Gate.SetStorageHandler(gate) //设置持久化处理器
}

/**
存储用户的Session信息
Session Bind Userid以后每次设置 settings都会调用一次Storage
*/
func (gate *Gate) Storage(Userid string, settings map[string]string) (err error) {
	log.Info("处理对Session的持久化, userid is %s, %v", Userid, settings)
	gate.storage.Set(Userid, settings)
	return nil
}

/**
强制删除Session信息
*/
func (gate *Gate) Delete(Userid string) (err error) {
	log.Info("删除Session持久化数据")
	gate.storage.Delete(Userid)
	return nil
}

/**
获取用户Session信息
用户登录以后会调用Query获取最新信息
*/
func (gate *Gate) Query(Userid string) (settings map[string]string, err error) {
	log.Info("查询Session持久化数据")


	if result := gate.storage.Get(Userid); result != nil {
		//settings = result
		return result, nil
	}
	return nil, fmt.Errorf("can not find the storage [%s] settings ", Userid)
}

/**
用户心跳,一般用户在线时60s发送一次
可以用来延长Session信息过期时间
*/
func (gate *Gate) Heartbeat(Userid string) {
	defer func() {
		log.Debug("更新心跳包， Userid is [%s]", Userid)
		if r := recover(); r != nil {
			log.Error("Gate Storage HeartBeat Error : %s", r)
		}
	}()

	if result := gate.storage.Get(Userid); result != nil {

		if timeout, ok := result["timeout"]; ok {
			if i, err := strconv.Atoi(timeout); err == nil && i < 120 {
				log.Info("更新用户 %s 在线的心跳包, now is %d", Userid, i + 60)
				result["timeout"] = string(i + 60)
			}
		}
	} else {
		log.Error("can not find the storage user :  [%s] ", Userid)
	}

}
