package apiV1

import (
  "gitlab.jiangxingai.com/asp-hrm/database"
  "gitlab.jiangxingai.com/asp-hrm/interfaces"
  "go.mongodb.org/mongo-driver/bson"
  log "k8s.io/klog"
)

//查询设备信息
func QueryDeviceMsg(deviceId string) (interfaces.CameraInfo, error) {
  var data interfaces.CameraInfo
  err := database.Db.QueryOne("device", bson.M{"_id": deviceId}).Decode(&data)
  if err != nil {
    return interfaces.CameraInfo{}, err
  }
  return data, nil
}

//插入日志
func InsertLog(data interfaces.Log) error {
  err := database.Db.Insert("log", &data)
  if err != nil {
    return err
  }
  return nil
}

//查询人员信息
func QueryHumanTitle(id string) (interfaces.Human, error) {
  var data interfaces.Human
  err := database.Db.QueryOne("human", bson.M{"_id": id}).Decode(&data)
  if err != nil {
    log.Error(err)
    return interfaces.Human{}, err
  }
  return data, nil
}

//更新或添加设备
func UpdateDeviceStatus(data interfaces.UpdateCameraInfo) error {
  err := database.Db.Update("device", bson.M{"_id": data.DeviceUuid}, bson.M{"$set": bson.M{"enable": data.Enable}}, true)
  if err != nil {
    return err
  }
  return nil
}

//添加设备
func InsertDevice(data interfaces.CameraInfo) error {
  err := database.Db.Insert("device", &data)
  if err != nil {
    return err
  }
  return nil
}

//删除设备
func DeleteDevice(id string) error {
  err := database.Db.Delete("device", bson.M{"_id": id})
  if err != nil {
    return err
  }
  return nil
}