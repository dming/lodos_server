package user

import (
	"github.com/dming/lodos/gate"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/dming/lodos/log"
	"github.com/dming/lodos/db/base"
)

func (m *user) Register(session gate.Session, args map[string]interface{}) error {
	mongoInfo := m.GetModuleSetting().Mongo
	coll, err := basedb.GetMongoFactories().GetCOLL(mongoInfo.Uri, mongoInfo.DB, USERCOLL)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	info := new(UserInfo)
	info.GetUserInfoFromMap(args)
	result := struct {
		A,B int
	}{}
	err = coll.Find(bson.M{"username": info.Username}).One(&result)
	if err == nil {
		log.Error("Username:%s alreay exist.", info.Username)
		return fmt.Errorf("Username:%s alreay exist.", info.Username)
	}

	err = coll.Insert(&UserInfo{
		Username: info.Username,
		Password: info.Password,
		Email:    info.Email,
	})
	if err != nil {
		log.Error("registe user of Username:%s failed.", info.Username)
		return fmt.Errorf("registe user of Username:%s failed.", info.Username)
	}
	log.Info("user success! %s, %s, %s", info.Username, info.Password, info.Email)
	return nil
}
