package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/astaxie/beego"
	//"errors"
	"fmt"
	"time"
)

type SystemConfig struct
{
	Id_  bson.ObjectId        `bson:"_id"`
	Name string               `bson:"name"`
	Nval int                  `bson:"nval"`
	Sval string               `bson:"sval"`
	Updated_at time.Time      `bson:"updated_at"`
	Created_at time.Time      `bson:"created_at"`
}


func GetAutoIncreaseId(db *mgo.Database, field_name string) (int,error)  {

	// get collection
	c := db.C("system_configs")
	var systemConfig SystemConfig
	err := c.Find(bson.M{"name": field_name}).One(&systemConfig)
	if err == mgo.ErrNotFound {
		currentTime := bson.Now()
		sc := SystemConfig{}
		sc.Id_ =  bson.NewObjectId()
		sc.Name = field_name
		sc.Nval = 1
		sc.Created_at = currentTime
		sc.Updated_at = currentTime
		insert_err := c.Insert(&sc)
		if insert_err !=nil{
			fmt.Println("insert system_config record failed")
			return -1,insert_err
		}
		return 1,nil
	} else if err != nil {
		return -1,err
	} else {
		change := mgo.Change{
			Update: bson.M{"$inc": bson.M{"nval": 1}},
			Upsert: true,
			ReturnNew: true,
		}

		_,err = c.Find(bson.M{"name": field_name}).Apply(change, &systemConfig)
		if err!=nil {
			return -1,err
		}

		return systemConfig.Nval,nil
	}

	return systemConfig.Nval,nil
}


func GetAccountDefaultCircles(db *mgo.Database) ([]SystemConfig,error)  {

	// get collection
	c := db.C("system_configs")
	var systemConfigList []SystemConfig
	err := c.Find(bson.M{"name": "circle_to_join","selector":"account"}).All(&systemConfigList)
	if err != nil {
		return nil,err
	} else {
		return systemConfigList,nil
	}

	return systemConfigList,nil
}

func FindSystemConfigByNameAndSelector(db *mgo.Database, name , selector string)(systemConfig SystemConfig, err error) {
	c := db.C("system_configs")
	err = c.Find(bson.M{"name": name,"selector":selector}).One(&systemConfig)
	return
}

func FindSystemConfigsByNameAndSelector(db *mgo.Database, name , selector string)(systemConfigs []SystemConfig, err error) {
	c := db.C("system_configs")
	err = c.Find(bson.M{"name": name,"selector":selector}).All(&systemConfigs)
	return
}