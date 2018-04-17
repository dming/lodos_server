package user

import (
	"github.com/dming/lodos/module/base"
	"github.com/globalsign/mgo"
	"github.com/dming/lodos/module"
	"github.com/dming/lodos/conf"
	"github.com/dming/lodos/log"
	"fmt"
	"github.com/dming/lodos/db/base"
)

const (
	COLL string = "CommonUsers"
)

var DB string

type UserInfo struct {
	Username string "bson:`username`"
	Password string "bson:`password`"
	Email    string "bson:`email`"
}

func (m *UserInfo) GetUserInfoFromMap(args map[string]interface{}) error {
	//info := new(UserInfo)
	if args["username"] == nil || args["password"] == nil {
		return fmt.Errorf("Username or Password cannot be nil")
	}
	m.Username = args["username"].(string)
	m.Password = args["password"].(string)
	if args["email"] != nil {
		m.Email = args["email"].(string)
	}
	log.Debug("%s, %s, %s", m.Username, m.Password, m.Email)
	return nil
}

var Module = func() module.Module {
	m := new(user)
	return m
}

type user struct {
	basemodule.BaseModule
	mongo *mgo.Session
}

func(m *user) GetType() string {
	return "User"
}

func(m *user) Version() string {
	return "v1.0.0"
}

func (m *user) OnInit(app module.AppInterface, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(app, m, settings)

	var err error
	m.mongo, err = basedb.GetMongoFactories().GetSession(settings.Mongo.Uri)
	if err != nil {
		log.Error(err.Error())
		m.mongo = nil
	}

	m.GetServer().RegisterGo("HD_Register", m.Register)
	m.GetServer().RegisterGo("HD_Login", m.login)
	m.GetServer().RegisterGo("Authentication", m.authentication)
}

func (m *user) Run(closeSig chan bool) {
	//TODO:
}

func (m *user) OnDestroy() {
	m.BaseModule.OnDestroy()
}
