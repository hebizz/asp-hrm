package interfaces

type CameraInfo struct {
	Id         string `json:"id" bson:"_id"`            //摄像头id
	Enable     int    `json:"enable" bson:"enable"`     //开关
	Uuid       string `json:"uuid" bson:"uuid"`         //host_id
	Interval   int    `json:"interval" bson:"interval"` //采集间隔
	DeviceName string `json:"name" bson:"name"`         //设备名称
	Location   string `json:"location" bson:"location"` //设备位置
}

type UpdateCameraInfo struct {
	DeviceUuid string `json:"device_uuid" bson:"device_uuid"` //摄像头id
	Enable     int    `json:"enable" bson:"enable"`           //开关
	Official   int    `json:"official" bson:"official"`       //人脸对比告警标识
}

type DeleteCameraInfo struct {
	DeviceId string `json:"device_id" bson:"device_id"` //摄像头id
}

type CaptureInfo struct {
	MsgId    string `json:"msg_id" bson:"msg_id"`       //msg_id
	DeviceId string `json:"device_id" bson:"device_id"` //摄像头id
	Base64   string `json:"img_str" bson:"base64"`      //base64
}

type Analysis struct {
	DeviceId  string          `json:"u" bson:"deviceId"`  //摄像头id
	Method    string          `json:"m" bson:"method"`    //auto / manual
	Timestamp int64           `json:"t" bson:"timestamp"` //时间戳
	Id        int64           `json:"id" bson:"image_id"` //img_id
	Position  [][]interface{} `json:"r" bson:"position"`  // [["ASP_HRM_FACE", "0.926804", 438, 245, 544, 403, "未知/李莉莉"], ...]
}

type PanguAlert struct {
	DeviceId string          `json:"device_id" bson:"deviceId"`      //摄像头id
	Title    string          `json:"title" bson:"title"`             //告警类型
	AlertMsg []AlertPosition `json:"alert_position" bson:"alertMsg"` //坐标
	Level    string          `json:"level" bson:"level"`             //告警等级
	ImageId  int64           `json:"image_Id" bson:"image_id"`       //img_id
}

type AlertPosition struct {
	Xmax        string `json:"xmax" bson:"xmax"`               //坐标
	Xmin        string `json:"xmin" bson:"xmin"`               //坐标
	Ymax        string `json:"ymax" bson:"ymax"`               //坐标
	Ymin        string `json:"ymin" bson:"ymin"`               //坐标
	Extra       string `json:"extra" bson:"extra"`             // 人员信息
	Probability string `json:"probability" bson:"probability"` //置信度
}

type UM301Msg struct {
	AudioType int `json:"audioType" bson:"audioType"` //人员类型
}
