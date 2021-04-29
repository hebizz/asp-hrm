package apiV1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.jiangxingai.com/asp-hrm/cloud"
	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	"gitlab.jiangxingai.com/asp-hrm/pkg/sdk"
	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "k8s.io/klog"
)

//向w2s发送获取图片topic
func TaskCaptureRoutine(camera interfaces.CameraInfo, client mqtt.Client) {
	log.Infof("start %s capture task", camera.Id)
	topic := fmt.Sprintf("%s/latest", camera.Uuid)
	payload := map[string]string{
		"msg_id":    "",
		"device_id": camera.Id,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return
	}
	for true {
		if token := client.Publish(topic, 2, false, bytes); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
			continue
		}
		//防止死循环loop
		if camera.Interval == 0 {
			camera.Interval = 2
		}
		time.Sleep(time.Duration(camera.Interval) * time.Second)
		res, err := QueryDeviceMsg(camera.Id)
		if err != nil {
			log.Errorf("can not find camera %s :: %s", camera.Id, err)
			break
		}
		if res.Enable == 0 {
			break
		}
	}
	log.Infof("disable %s capture task", camera.Id)
}

//从w2s接收图片, sdk识别
var CaptureAnalysisHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var data interfaces.CaptureInfo
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Error("unmarshal edge info error: ", err)
		return
	}
	sdkRes, flag, err := sdk.AnalysisFaceBaiduSdk(data.Base64, data.DeviceId)
	if err != nil {
		return
	}
	//识别到人脸
	if flag {
		go PubAnalysisDataToPangu(client, data, sdkRes)
	}
	return
}

//获取pangu告警结果
var GetAiResultHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var data interfaces.PanguAlert
	var path, aiPath string
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Error("unmarshal pangu info error: ", err)
		return
	}
	log.Infof("receive pangu alert data:: %+v", data)
	imgStr, ok := sdk.ImageMap.Load(data.ImageId)
	if !ok {
		log.Error("can not find imageId")
		return
	}
	aiPath, err := utils.Base64ToFile(data.ImageId, imgStr.(string))
	if err != nil {
		log.Error(err)
		return
	}
	sdk.ImageMap.Delete(data.ImageId)
	deviceMsg, err := QueryDeviceMsg(data.DeviceId)
	if err != nil {
		log.Error(err)
		return
	}
	for _, v := range data.AlertMsg {
		var humanType, title, Dname string
		if v.Extra == "未知" {
			humanType = "1"
			title = v.Extra
			Dname = "-"
			path = ""
		} else {
			humanType = "0"
			res, _ := QueryHumanTitle(v.Extra)
			title = res.Title
			Dname = res.Dname
			path = res.Path
		}
		aiLog := interfaces.Log{
			Id:        primitive.NewObjectID().Hex(),
			HumanType: humanType,
			AiResult:  v,
			TimeStamp: time.Now().Unix(),
			Title:     title,
			Device:    deviceMsg.DeviceName,
			Position:  deviceMsg.Location,
			Dname:     Dname,
			RawPath:   path,
			AiPath:    fmt.Sprintf("/api/%s", strings.Split(aiPath, "/data/local/asp/")[1]),
		}
		if err := InsertLog(aiLog); err != nil {
			log.Error(err)
			return
		}
		go PubHumanTypeToUM301(client, data.DeviceId, humanType)
	}
	return
}

//更新或者添加摄像头
var UpdateCameraHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var data interfaces.UpdateCameraInfo
	var cameraList []interfaces.CameraInfo
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Error("unmarshal pangu update camera info error: ", err)
		return
	}
	log.Infof("receive pangu update camera data:: %+v", data)
	res, flag := QueryDeviceMsg(data.DeviceUuid)
	if flag != nil {
		log.Errorf("add camera %s :: %s", data.DeviceUuid, flag)
		cameraList = append(cameraList, interfaces.CameraInfo{Id: data.DeviceUuid, Enable: data.Enable})
		cameraList, err := cloud.GetCamerasFromAsp(cameraList)
		if err != nil {
			log.Error(err)
			return
		}
		if len(cameraList) > 0 {
			err = InsertDevice(cameraList[0])
			if err != nil {
				log.Error(err)
				return
			}
			if cameraList[0].Enable == 1 && data.Official == 2{
				go TaskCaptureRoutine(cameraList[0], client)
				return
			}
		}
	} else {
		err := UpdateDeviceStatus(data)
		if err != nil {
			log.Error(err)
			return
		}
	}
	if res.Enable != data.Enable && data.Enable == 1 {
		go TaskCaptureRoutine(res, client)
	}
	return
}

//删除摄像头
var DeleteCameraHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var data interfaces.DeleteCameraInfo
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Error("unmarshal pangu delete camera info error: ", err)
		return
	}
	log.Infof("receive pangu delete camera data:: %+v", data)
	if err := DeleteDevice(data.DeviceId); err != nil {
		log.Error(err)
		return
	}
}

//发送识别信息给pangu
func PubAnalysisDataToPangu(client mqtt.Client, data interfaces.CaptureInfo, sdkRes interfaces.Analysis) {
	topic := fmt.Sprintf("pangu/alert/%s", data.DeviceId)
	bytes, err := json.Marshal(sdkRes)
	if err != nil {
		log.Error(err)
		return
	}
	if token := client.Publish(topic, 2, false, bytes); token.Error() != nil {
		log.Error(token.Error())
		return
	}
	log.Infof("send data to pangu:: %+v", sdkRes)
}

//发送人员类型告警给UM301
func PubHumanTypeToUM301(client mqtt.Client, deviceId string, humanType string) {
	audioType, err := strconv.Atoi(humanType)
	if err != nil {
		log.Error(err)
		return
	}
	topic := fmt.Sprintf("%s/hrm/audio", deviceId)
	bytes, err := json.Marshal(interfaces.UM301Msg{AudioType: audioType})
	if err != nil {
		log.Error(err)
		return
	}
	if token := client.Publish(topic, 2, false, bytes); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return
	}
}
