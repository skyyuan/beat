package utils

import (
	"gopkg.in/mgo.v2"
)


var session *mgo.Session


func init() {
	InitMgo()
}


func ConnMgo() *mgo.Session {
	return session.Copy()
}

func InitMgo() {
	sess, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	session = sess.Copy()

	session.SetMode(mgo.Monotonic, true)
}

func GetMgoDbSession() (*mgo.Database, *mgo.Session){
	mgoSession := ConnMgo()
	db := mgoSession.DB("skiing")
	return db, mgoSession
}