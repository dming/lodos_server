package main

import (
	"github.com/dming/lodos"
	"server/gate"
	"server/login"
	LodosGate "github.com/dming/lodos/gate"
	"github.com/dming/lodos/module"
	"server/user"
)
//func ChatRoute( app module.app,Type string,hash string) (*module.ServerSession){
//	//演示多个服务路由 默认使用第一个Server
//	log.Debug("Hash:%s 将要调用 type : %s",hash,Type)
//	servers:=app.GetServersByType(Type)
//	if len(servers)==0{
//		return nil
//	}
//	return servers[0]
//}
var app module.AppInterface

func main() {
	app = lodos.CreateApp("1.0.0")
	//app.Route("Chat",ChatRoute)
	app.SetJudgeGuest(judgeGuest)
	app.Run(true, //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		gate.Module(),  //这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		login.Module(), //这是用户登录验证模块
		//chat.Module(),
		//webapp.Module(),
		user.Module(),
	)
}

//true is guest, false is user
func judgeGuest(session LodosGate.Session) bool {
	username := session.GetUserid()
	if username == "" {
		return true
	}
	token := session.GetToken()
	if token == "" {
		return true
	}
	serverId := session.GetServerid()
	if serverId == "" {
		return true
	}
	server, err := app.GetServerById(serverId)
	if err != nil {
		return true
	}
	results, err := server.Call("Authentication", username, token)
	if err != nil {
		return true
	}
	if results != nil && len(results) > 0 {
		if re, ok := results[0].(bool); ok {
			return !re //authentication is user, otherwise is guest
		}
	}
	return true
}