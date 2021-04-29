package sdk

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"sync"
	"time"

	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	e "gitlab.jiangxingai.com/asp-hrm/pkg/response"
	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	log "k8s.io/klog"
)

var BaiduSdkGroupId = utils.GetEnv("BAIDUSDK_GROUPID", "1")

const (
	BaiduSdkGrantType    = "client_credentials"
	BaiduSdkClientId     = "SrYzl1TqukNY9Nb6I5o4w3OH"
	BaiduSdkClientSecret = "Edfw9F4AlwYV73fyBGHB958ph3R6oE26"
	BaiduSdkTokenUrl     = "https://aip.baidubce.com/oauth/2.0/token"
	BaiduSdkRegisterUrl  = "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add"
	BaiduSdkUpdateUrl    = "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/update"
	BaiduSdkRemoveUrl    = "https://aip.baidubce.com/rest/2.0/face/v3/faceset/face/delete"
	BaiduSdkAnalysisUrl  = "https://aip.baidubce.com/rest/2.0/face/v3/multi-search"
	BaiduSdkDetectUrl    = "https://aip.baidubce.com/rest/2.0/face/v3/detect"
)

//如果并发量大, 考虑map加锁, 写效率更高
var ImageMap sync.Map

//var ImageStruct = struct {
//	sync.RWMutex
//	ImageMap map[int64]string
//}{
//	ImageMap: make(map[int64]string, 100),
//}

//获取token
func GetBaiduSdkToken() (string, error) {
	params := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s",
		url.QueryEscape(BaiduSdkGrantType), url.QueryEscape(BaiduSdkClientId),
		url.QueryEscape(BaiduSdkClientSecret))
	tokenUrl := fmt.Sprintf("%s?%s", BaiduSdkTokenUrl, params)
	body, err := utils.HttpGet(tokenUrl)
	if err != nil {
		return "", err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return "", err
	}
	return res["access_token"].(string), nil
}

//注册人员, 更新人员到百度人脸库
func RegisterHumanToBaiduSdk(feature string, data interfaces.Human) (string, error) {
	var featureUrl, featureData string
	token, err := GetBaiduSdkToken()
	if err != nil {
		return "", err
	}
	fileData, err := ioutil.ReadFile(data.Rpath)
	if err != nil {
		return "", err
	}
	if feature == "register" {
		featureUrl = fmt.Sprintf("%s?access_token=%s", BaiduSdkRegisterUrl, token)
		featureData = fmt.Sprintf("image=%s&image_type=BASE64&group_id=%s&user_id=%s&action_type=REPLACE",
			url.QueryEscape(base64.StdEncoding.EncodeToString(fileData)), BaiduSdkGroupId, data.Id)
	} else {
		featureUrl = fmt.Sprintf("%s?access_token=%s", BaiduSdkUpdateUrl, token)
		featureData = fmt.Sprintf("image=%s&image_type=BASE64&group_id=%s&user_id=%s&action_type=UPDATE",
			url.QueryEscape(base64.StdEncoding.EncodeToString(fileData)), BaiduSdkGroupId, data.Id)
	}

	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		return "", err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return "", err
	}
	log.Info("sdk register human res: %+v", res)
	if res["error_code"].(float64) != 0 {
		return "", errors.New(res["error_msg"].(string))
	}
	return res["result"].(map[string]interface{})["face_token"].(string), nil
}

//删除人员
func RemoveHumanToBaiduSdk(data interfaces.Human) error {
	token, err := GetBaiduSdkToken()
	if err != nil {
		return err
	}
	featureUrl := fmt.Sprintf("%s?access_token=%s", BaiduSdkRemoveUrl, token)
	featureData := fmt.Sprintf("group_id=%s&user_id=%s&face_token=%s", BaiduSdkGroupId, data.Id, data.FaceToken)
	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		return err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return err
	}
	log.Infof("sdk remove human res: %+v", res)
	if res["error_code"].(float64) != 0 {
		return errors.New(res["error_msg"].(string))
	}
	return nil
}

//人脸搜索
func AnalysisFaceBaiduSdk(baseStr string, deviceId string) (interfaces.Analysis, bool, error) {
	data := interfaces.Analysis{
		DeviceId:  deviceId,
		Method:    "a",
		Timestamp: time.Now().Unix(),
		Id:        time.Now().Unix(),
	}
	flag := false
	token, err := GetBaiduSdkToken()
	if err != nil {
		log.Error(err)
		return interfaces.Analysis{}, false, err
	}
	featureUrl := fmt.Sprintf("%s?access_token=%s", BaiduSdkAnalysisUrl, token)
	featureData := fmt.Sprintf("image=%s&image_type=BASE64&group_id_list=%s&max_face_num=10&max_user_num=10&quality_control=LOW",
		url.QueryEscape(baseStr), BaiduSdkGroupId)
	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		log.Error(err)
		return interfaces.Analysis{}, false, err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		log.Error(err)
		return interfaces.Analysis{}, false, err
	}
	log.Infof("device %s execute sdk analysis result:: %v", deviceId, res)
	if res["error_code"].(float64) == 0 {
		faceList := res["result"].(map[string]interface{})["face_list"].([]interface{})
		data.Position = GenerateAnalysis(faceList)
	}
	// 人脸搜索未识别到人, 调用人脸检测
	if res["error_code"].(float64) == 222207 {
		faceList, err := AnalysisOtherHumanOnBaiduSdk(baseStr, deviceId)
		if err != nil {
			return interfaces.Analysis{}, false, err
		}
		data.Position = GenerateAnalysis(faceList)
	}
	// 如果有识别结果, 放入队列
	if len(data.Position) > 0 {
		count := 0
		ImageMap.Range(func(k, v interface{}) bool {
			count++
			return true
		})
		log.Info("ImageMap len is::", count)
		if count == 100 {
			log.Info("ImageMap is full")
			ImageMap.Range(func(k, v interface{}) bool {
				//清除10s之前缓存
				if data.Timestamp-k.(int64) > 10 {
					ImageMap.Delete(k)
				}
				return true
			})
		}
		ImageMap.Store(data.Timestamp, baseStr)
		flag = true
	}
	return data, flag, nil
}

//生成pngu所需的数据结构
func GenerateAnalysis(faceList []interface{}) [][]interface{} {
	var position [][]interface{}
	for _, v := range faceList {
		var aiResult []interface{}
		aiResult = append(aiResult, "ASP_HRM_FACE")
		_, ok := v.(map[string]interface{})["user_list"]
		if ok {
			if len(v.(map[string]interface{})["user_list"].([]interface{})) > 0 {
				aiResult = append(aiResult, fmt.Sprint(v.(map[string]interface{})["user_list"].([]interface{})[0].(map[string]interface{})["score"].(float64)/100))
			} else {
				aiResult = append(aiResult, "1")
			}
		} else {
			aiResult = append(aiResult, "1")
		}
		aiResult = append(aiResult, v.(map[string]interface{})["location"].(map[string]interface{})["left"].(float64))
		aiResult = append(aiResult, v.(map[string]interface{})["location"].(map[string]interface{})["top"].(float64))
		aiResult = append(aiResult, v.(map[string]interface{})["location"].(map[string]interface{})["width"].(float64)+aiResult[2].(float64))
		aiResult = append(aiResult, v.(map[string]interface{})["location"].(map[string]interface{})["height"].(float64)+aiResult[3].(float64))
		if aiResult[1] != "1" {
			aiResult = append(aiResult, v.(map[string]interface{})["user_list"].([]interface{})[0].(map[string]interface{})["user_id"])
		} else {
			aiResult = append(aiResult, "未知")
		}
		position = append(position, aiResult)
	}
	return position
}

//人脸对比, 判断人脸是否录入库中
func CheckHumanOnBaiduSdk(path string) (int, error) {
	token, err := GetBaiduSdkToken()
	if err != nil {
		return 0, err
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	featureUrl := fmt.Sprintf("%s?access_token=%s", BaiduSdkAnalysisUrl, token)
	featureData := fmt.Sprintf("image=%s&image_type=BASE64&group_id_list=%s&max_face_num=10&max_user_num=10&quality_control=LOW",
		url.QueryEscape(base64.StdEncoding.EncodeToString(fileData)), BaiduSdkGroupId)
	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		return 0, err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return 0, err
	}
	log.Infof("sdk search res: %+v", res)
	if res["error_code"].(float64) == 0 {
		return e.ERROR_SAME_HUMAN, nil
	}
	return 0, nil
}

//人脸检测, 检测是否为人脸, 是否存在多张人脸
func CheckMultiHumanOnBaiduSdk(path string) (int, error) {
	token, err := GetBaiduSdkToken()
	if err != nil {
		return 0, err
	}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	featureUrl := fmt.Sprintf("%s?access_token=%s", BaiduSdkDetectUrl, token)
	featureData := fmt.Sprintf("image=%s&image_type=BASE64&max_face_num=10&liveness_control=NORMAL",
		url.QueryEscape(base64.StdEncoding.EncodeToString(fileData)))
	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		return 0, err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return 0, err
	}
	log.Infof("sdk detect res: %+v", res)
	if res["error_code"].(float64) == 222202 {
		return e.ERROR_QUALITY_HUMAN, nil
	}
	if res["result"].(map[string]interface{})["face_num"].(float64) > 1 {
		return e.ERROR_MUTLI_HUMAN, nil
	}
	return 0, nil
}

//人脸检测, 识别系统中陌生人
func AnalysisOtherHumanOnBaiduSdk(baseStr string, deviceId string) ([]interface{}, error) {
	token, err := GetBaiduSdkToken()
	if err != nil {
		return nil, err
	}
	featureUrl := fmt.Sprintf("%s?access_token=%s", BaiduSdkDetectUrl, token)
	featureData := fmt.Sprintf("image=%s&image_type=BASE64&max_face_num=10&liveness_control=NORMAL",
		url.QueryEscape(baseStr))
	body, err := utils.HttpPost(featureUrl, featureData)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	}
	log.Infof("device %s sdk detect res: %+v", deviceId, res)
	if res["error_code"].(float64) != 0 {
		return nil, errors.New(res["error_msg"].(string))
	} else {
		return res["result"].(map[string]interface{})["face_list"].([]interface{}), nil
	}
}
