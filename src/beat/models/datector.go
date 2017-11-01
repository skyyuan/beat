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
	Type           string        `bson:type`
	Factor         float64       `bson:factor`
	CommonModel    `bson:",inline"`
}

func NewDetector(db *mgo.Database, DeviceId, tp, Ip string) (detector Detector, err error) {
	currentTime := bson.Now()
	detector.Id_ = bson.NewObjectId()
	detector.DetectorId, _ = GetAutoIncreaseId(db, "detector_id")
	detector.DeviceId = DeviceId
	detector.Ip = Ip
	detector.Status = "true"
	detector.Type = tp
	detector.CreatedAt = currentTime
	detector.UpdatedAt = currentTime
	adminuserCollection := db.C("detectors")
	err = adminuserCollection.Insert(&detector)
	return
}

func FindByDetector(db *mgo.Database, id int) (detector Detector, err error) {
	collection := db.C("detectors")
	err = collection.Find(bson.M{"detector_id": id}).One(&detector)
	return
}

func GetDetectors(db *mgo.Database) (detectors []Detector, err error) {
	collection := db.C("detectors")
	err = collection.Find(nil).All(&detectors)
	return
}

func GetAllDetectors(db *mgo.Database, page, perPage int, deviceId string) (detectors []Detector, err error) {
	collection := db.C("detectors")
	query := bson.M{}
	if deviceId != "" {
		query["device_id"] = bson.M{"$regex": deviceId}
	}
	err = collection.Find(query).Limit(perPage).Skip((page - 1) * perPage).Sort("-created_at").All(&detectors)
	return
}

func GetAllDetectorsCount(db *mgo.Database, deviceId string) (count int, err error) {
	collection := db.C("detectors")
	query := bson.M{}
	if deviceId != "" {
		query["device_id"] = bson.M{"$regex": deviceId}
	}
	count, err = collection.Find(query).Count()
	return
}

func GetDetectorByDeviceId(db *mgo.Database, deviceId string) (detector Detector, err error) {
	collection := db.C("detectors")
	err = collection.Find(bson.M{"device_id": deviceId}).One(&detector)
	return
}

func DeleteDetector(db *mgo.Database, id int) (err error) {
	collection := db.C("detectors")
	err = collection.Remove(bson.M{"detector_id": id})
	return
}

func (self *Detector) UpdateByParams(db *mgo.Database, DeviceId, location, ip string, factor float64,  prevDetectorId, nextDetectorId string) (err error) {
	query := bson.M{"factor": factor,"device_id": DeviceId, "location": location, "ip": ip, "updated_at": bson.Now()}
	if prevDetectorId != "" {
		query["prev_detector_id"] = bson.ObjectIdHex(prevDetectorId)
	}else {
		query["prev_detector_id"] = nil
	}
	if nextDetectorId != "" {
		query["next_detector_id"] = bson.ObjectIdHex(nextDetectorId)
	}else {
		query["next_detector_id"] = nil
	}
	userCollection := db.C("detectors")
	err = userCollection.Update(bson.M{"_id": self.Id_}, bson.M{"$set": query})
	return
}
