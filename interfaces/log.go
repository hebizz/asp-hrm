package interfaces

type Log struct {
	Id        string        `json:"id" bson:"_id"`              // 主键id
	HumanType string        `json:"humanType" bson:"humanType"` // 人员类型{0:系统人员,1:陌生人}
	Title     string        `json:"title" bson:"title"`         // 人员姓名[..., 未知]
	Dname     string        `json:"dname" bson:"dname"`         // 所属部门[..., -]
	TimeStamp int64         `json:"timestamp" bson:"timestamp"` // 识别时间
	Device    string        `json:"device" bson:"device"`       // 设备名称
	Position  string        `json:"position" bson:"position"`   // 安装位置
	RawPath   string        `json:"rawPath" bson:"rawPath"`     // 原始照片
	AiPath    string        `json:"aiPath" bson:"aiPath"`       // ai识别照片
	AiResult  AlertPosition `json:"aiResult" bson:"aiResult"`   //坐标
}
