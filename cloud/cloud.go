package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"gitlab.jiangxingai.com/asp-hrm/database"
	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"

	log "k8s.io/klog"
)

const (
	PanguAlertParam = "ASP_HRM_FACE"
	PanguAlertUrl   = ":9005/pangu/api/v1/rule/list?type="
	AspDeviceUrl    = ":8085/api/v1/feign/devices/info/"
)

func Setup() []interfaces.CameraInfo {
	cameraList, err := GetCamerasFromPangu()
	for ; err != nil; {
		log.Error(err)
		time.Sleep(time.Duration(2) * time.Second)
		cameraList, err = GetCamerasFromPangu()
	}
	cameraList, err = GetCamerasFromAsp(cameraList)
	for ; err != nil; {
		log.Error(err)
		time.Sleep(time.Duration(2) * time.Second)
		cameraList, err = GetCamerasFromAsp(cameraList)
	}
	log.Info("get cloud cameraList: ", cameraList)
	for _, v := range cameraList {
		if err = database.Db.Update("device", bson.M{"_id": v.Id}, bson.M{"$set": &v}, true); err != nil {
			log.Error(err)
		}
	}
	return cameraList
}

//向pangu获取人脸设备摄像头id
func GetCamerasFromPangu() ([]interfaces.CameraInfo, error) {
	params := url.QueryEscape(PanguAlertParam)
	cloudAddr := utils.GetEnv("PANGU_ADDR", "10.56.0.52")
	if cloudAddr == "" {
		return nil, errors.New("pangu server addr can not find")
	}
	getCameraUrl := fmt.Sprintf("http://%s%s%s", cloudAddr, PanguAlertUrl, params)
	body, err := utils.HttpGet(getCameraUrl)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	}
	log.Info(res)
	if res["data"].(map[string]interface{})["rules"] == nil {
		return nil, err
	}
	var camera interfaces.CameraInfo
	var cameraList []interfaces.CameraInfo
	cameras := res["data"].(map[string]interface{})["rules"].([]interface{})
	for _, v := range cameras {
		camera.Id = v.(map[string]interface{})["device_uuid"].(string)
		camera.Enable = int(v.(map[string]interface{})["enable"].(float64))
		cameraList = append(cameraList, camera)
	}
	log.Info(cameraList)
	return cameraList, nil
}

//根据摄像头id向asp获取host_uuid, interval间隔
func GetCamerasFromAsp(cameraList []interfaces.CameraInfo) ([]interfaces.CameraInfo, error) {
	for index, camera := range cameraList {
		params := camera.Id
		cloudAddr := utils.GetEnv("ASP_ADDR", "10.56.0.52")
		if cloudAddr == "" {
			return cameraList, errors.New("pangu server addr can not find")
		}
		getCameraUrl := fmt.Sprintf("http://%s%s%s", cloudAddr, AspDeviceUrl, params)
		log.Info(getCameraUrl)
		body, err := utils.HttpGet(getCameraUrl)
		if err != nil {
			return cameraList, err
		}
		res := make(map[string]interface{})
		if err = json.Unmarshal([]byte(body), &res); err != nil {
			return cameraList, err
		}
		log.Info(res)
		//host_uuid为空, 此设备为智慧眼设备
		if res["data"] == nil {
			return cameraList, errors.New("can not get camera msg from asp")
		}
		if res["data"].(map[string]interface{})["device"].(map[string]interface{})["host_uuid"] == nil {
			cameraList[index].Uuid = cameraList[index].Id
		} else if res["data"].(map[string]interface{})["device"].(map[string]interface{})["host_uuid"].(string) == "" {
			cameraList[index].Uuid = cameraList[index].Id
		} else {
			cameraList[index].Uuid = res["data"].(map[string]interface{})["device"].(map[string]interface{})["host_uuid"].(string)
		}
		cameraList[index].Interval = int(res["data"].(map[string]interface{})["setting"].(map[string]interface{})["interval"].(float64))
		cameraList[index].DeviceName = res["data"].(map[string]interface{})["device"].(map[string]interface{})["name"].(string)
		cameraList[index].Location = res["data"].(map[string]interface{})["device"].(map[string]interface{})["location"].(string)
	}
	return cameraList, nil
}
