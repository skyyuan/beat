package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//探测器表
type Detector struct {
	Id_            bson.ObjectId `bson:"_id"`
	DetectorId     int           `bson:"detector_id"`
	DeviceId       string        `bson:"device_id"`
	PrevDetectorId bson.ObjectId `bson:"prev_detector_id,omitempty"`
	NextDetectorId bson.ObjectId `bson:"next_detector_id,omitempty"`
	Location       string        `bson:"location"`
	Ip             string        `bson:ip`
	Status         string        `bson:status`
	CheckStatus    string        `bson:check_status`
	Type           string        `bson:type`
	Factor         float64       `bson:factor`
	CommonModel    `bson:",inline"`
}

func NewDetector(db *mgo.Database, DeviceId, tp, Ip, checkStatus string) (detector Detector, err error) {
	currentTime := bson.Now()
	detector.Id_ = bson.NewObjectId()
	detector.DetectorId, _ = GetAutoIncreaseId(db, "detector_id")
	detector.DeviceId = DeviceId
	detector.Ip = Ip
	detector.Status = "true"
	detector.CheckStatus = checkStatus
	detector.Type = tp
	detector.CreatedAt = currentTime
	detector.UpdatedAt = currentTime
	adminuserCollection := db.C("detectors")
	err = adminuserCollection.Insert(&detector)
	return
}

func GetDetectorByDeviceId(db *mgo.Database, deviceId, tp string) (detector Detector, err error) {
	collection := db.C("detectors")
	err = collection.Find(bson.M{"device_id": deviceId, "type": tp}).One(&detector)
	return
}

func (self *Detector) UpdateByParams(db *mgo.Database, ip, checkStatus string) (err error) {
	query := bson.M{"ip": ip, "check_status": checkStatus,  "updated_at": bson.Now()}
	userCollection := db.C("detectors")
	err = userCollection.Update(bson.M{"_id": self.Id_}, bson.M{"$set": query})
	return
}

func (self *Detector) UpdateByStatus(db *mgo.Database) (err error) {
	query := bson.M{"status": "true", "updated_at": bson.Now()}
	userCollection := db.C("detectors")
	err = userCollection.Update(bson.M{"_id": self.Id_}, bson.M{"$set": query})
	return
}
