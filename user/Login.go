package user

import (
	"github.com/dming/lodos/gate"
	"fmt"
	"github.com/dming/lodos/db/base"
	"github.com/globalsign/mgo/bson"
	"github.com/dming/lodos/gate/base"
	"github.com/dming/lodos/log"
)

// return token ?
func (m *user) login(session gate.Session, msg map[string]interface{}) (token string, err error) {
	//time.Sleep(time.Millisecond * 200)
	if msg["username"] == nil || msg["password"] == nil {
		err = fmt.Errorf("userName or passWord cannot be nil")
		return "", err
	}
	var username string
	if u, ok := msg["username"].(string); !ok {
		err = fmt.Errorf("username cannot be format as [string]")
		return "", err
	} else {
		username = u
	}
	var password string
	if p, ok := msg["password"].(string); !ok {
		err = fmt.Errorf("password cannot be format as [string]")
		return "", err
	} else {
		password = p
	}

	//TODO: check if valid of username and password, if valid, set token, else token remain to ""
	mongoInfo := m.GetModuleSetting().Mongo
	coll, err := basedb.GetMongoFactories().GetCOLL(mongoInfo.Uri, mongoInfo.DB, COLL)
	if err != nil {
		return "", err
	}
	result := struct {
		A,B int
	}{}
	err = coll.Find(bson.M{"username": username, "password": password}).One(&result)
	if err != nil {
		return "", err
	} else {
		token, err = m.createToken(username)
		if err != nil {
			return "", err
		}
	}

	if token == "" {
		return "", fmt.Errorf("login fail with %s and %s", username, password)
	}

	session.SetToken(token)
	err = session.Bind(username)
	if err != nil {
		return token, err
	}
	session.SetToken(token)
	err = session.Push() //推送到网关
	if err != nil {
		log.Error("cannot push , err is %s", err.Error())
		return token, err
	}
	return token, nil
}

func (m *user) authentication(username string, token string) (bool, error) {
	//TODO: authenticate the username and token in the database
	return false, nil
}

func (m *user) createToken(username string) (string, error) {
	//todo: create token by given username
	return basegate.Get_uuid(), nil
}
