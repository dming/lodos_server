/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package gate

import (
	"github.com/dming/lodos/conf"
	"github.com/dming/lodos/module"
	"github.com/dming/lodos/gate/base"
)

var Module = func() module.Module {
	gate := new(Gate)
	return gate
}

type Gate struct {
	basegate.Gate //继承

	//storage map[string]map[string]string
	//storage *safemap.BeeMap4MapStr
}

func (gate *Gate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Gate"
}
func (gate *Gate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (gate *Gate) OnInit(app module.AppInterface, settings *conf.ModuleSettings) {
	//注意这里一定要用 gate.Gate 而不是 module.BaseModule
	gate.Gate.OnInit(app, gate, settings)

	//gate.storage = safemap.NewBeeMap4MapStr()
	gate.Gate.SetStorageHandler(gate) //设置持久化处理器
}
func (gate *Gate) Run(closeSig chan bool) {
	gate.Gate.Run(closeSig)
}
func (gate *Gate) OnDestroy() {
	gate.Gate.OnDestroy()
}
