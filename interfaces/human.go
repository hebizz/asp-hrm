package interfaces

type Department struct {
	Id        string `json:"id" bson:"_id"`              // 主键,部门id
	Pid       string `json:"pid" bson:"pid"`             // 上级部门id
	Title     string `json:"title" bson:"title"`         // 部门名称
	Type      string `json:"type" bson:"type"`           // 类型
	Timestamp int64  `json:"timestamp" bson:"timestamp"` //时间戳
}

type DepartmentTree struct {
	Id       string           `json:"id" bson:"_id"`            // 主键,部门id
	Pid      string           `json:"pid" bson:"pid"`           // 上级部门id
	Title    string           `json:"title" bson:"title"`       // 部门名称
	Children []DepartmentTree `json:"children" bson:"children"` //子部门
	HList    []Human          `json:"hList" bson:"hList"`       //子用户列表
	Type     string           `json:"type" bson:"type"`         //类型
}

type Human struct {
	Id        string `json:"id" bson:"_id"`        //主键id
	UserId    string `json:"userId" bson:"userId"` //用户id
	Pid       string `json:"pid" bson:"pid"`       //部门id
	Dname     string `json:"dname" bson:"dname"`   //部门名称
	Title     string `json:"title" bson:"title"`   //用户名字
	Path      string `json:"path" bson:"path"`     //头像路径
	Rpath     string `json:"rpath" bson:"rpath"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"` //录入时间
	Type      string `json:"type" bson:"type"`           //类型
	FaceToken string `json:"faceToken" bson:"faceToken"` //百度sdk人脸token
}

type HumanIds struct {
	Ids []string `json:"ids" binding:"required"` //用户id列表
}
