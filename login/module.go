/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package login

import (
	"fmt"
	"github.com/dming/lodos/conf"
	"github.com/dming/lodos/gate"
	"github.com/dming/lodos/module"
	"github.com/dming/lodos/module/base"
	"github.com/dming/lodos/log"
	"github.com/go-redis/redis"
	"github.com/dming/lodos/utils/uuid"
)

var Module = func() module.Module {
	m := new(Login)
	return m
}

type Login struct {
	basemodule.BaseModule
	db *redis.Client
}

func (m *Login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *Login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Login) OnInit(app module.AppInterface, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(app, m, settings)
	//m.SetDBClient("127.0.0.1:6379", "")
	//m.GetServer().RegisterGo("HD_Login", m.login)  //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGo("getRand", m.getRand) //演示后台模块间的rpc调用
	m.GetServer().Register("HD_Robot", m.robot)
	m.GetServer().RegisterGo("HD_Robot_GO", m.robot)  //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGo("HD_Test", m.Test)
}

func (m *Login) Run(closeSig chan bool) {
}

func (m *Login) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}
func (m *Login) robot(session gate.Session, msg map[string]interface{}) (result string, err error) {
	//time.Sleep(1)
	//log.Info("function on call robot:  %s", string(r))
	if msg["userName"] == nil || msg["passWord"] == nil {
		err = fmt.Errorf("userName or passWord cannot be nil")
		return
	}
	return fmt.Sprintf("%s, %s", msg["userName"], msg["passWord"]), nil
}


func (m *Login) getRand(session gate.Session, by []byte,mp map[string]interface{},f float64,i int32,b bool) (result string, err error) {
	//演示后台模块间的rpc调用
	return fmt.Sprintf("My is Login Module %s", by,mp,f,i,b), nil
}

func (m *Login) Test(session gate.Session, args map[string]interface{}) error {
	log.Info("IP is %s,\n serverID is %s,\n sessionID is %s,\n userID is %s, \n",
		session.GetIP(), session.GetServerid(), session.GetSessionid(), session.GetUserid())

	var msg string = ""
	if args["msg"] != nil {//&& reflect.TypeOf(args["msg"]) == reflect.TypeOf(reflect.String) {
		msg = args["msg"].(string)
	}
	uuid := uuid.Rand()
	log.Info("Get a test request from mqtt client,\n uuid: %s \n message is %s", uuid.Hex(), msg)
	return nil
}

func (m *Login) SetDBClient(addr string, password string, args ...interface{}) (error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr, //"127.0.0.1:6379",
		Password: password, // "guest",
		DB : 0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Error("connect to redis %s fail, err is  [%s]", addr, err)
		return err
	}
	m.db = client
	log.Info("connect to redis %s success, ping result [%s]", addr, pong)

	//
	return nil
}