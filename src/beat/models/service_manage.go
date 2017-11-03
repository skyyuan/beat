package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//服务管理表
type ServiceManage struct {
	Id_             bson.ObjectId `bson:"_id"`
	ServiceManageId int           `bson:"service_manage_id"`
	ServiceName     string        `bson:"service_name"`
	Status          string        `bson:status`
	CommonModel     `bson:",inline"`
}

func NewServiceManage(db *mgo.Database, name string) (service_name ServiceManage, err error) {
	currentTime := bson.Now()
	service_name.Id_ = bson.NewObjectId()
	service_name.ServiceManageId, _ = GetAutoIncreaseId(db, "service_manage_id")
	service_name.ServiceName = name
	service_name.Status = "true"
	service_name.CreatedAt = currentTime
	service_name.UpdatedAt = currentTime
	Collection := db.C("service_manages")
	err = Collection.Insert(&service_name)
	return
}

func GetServiceManage(db *mgo.Database, name string) (service_name ServiceManage, err error) {
	collection := db.C("service_manages")
	err = collection.Find(bson.M{"service_name": name}).One(&service_name)
	return
}

func (self *ServiceManage) UpdateByStatus(db *mgo.Database) (err error) {
	query := bson.M{"status": "true", "updated_at": bson.Now()}
	userCollection := db.C("service_manages")
	err = userCollection.Update(bson.M{"_id": self.Id_}, bson.M{"$set": query})
	return
}
