package apiV1

import (
	"fmt"
	"strconv"

	database "gitlab.jiangxingai.com/asp-hrm/database"
	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	log "k8s.io/klog"
)

//获取所有部门信息
func QueryDepartmentList(departmentList interface{}) (interface{}, error) {
	opts := options.FindOptions{
		Sort: bson.M{"timestamp": -1},
	}
	cursor, ctx, _ := database.Db.QueryAll("department", bson.M{}, &opts)
	defer cursor.Close(ctx)
	err := cursor.All(ctx, &departmentList)
	if err != nil {
		return nil, err
	}
	return departmentList, nil
}

//获取子部门信息
func GetChildrenDepartment(id string) ([]interfaces.Department, error) {
	var departmentList, childrenDepartmentList []interfaces.Department
	cursor, ctx, _ := database.Db.QueryAll("department", bson.M{"pid": id}, nil)
	defer cursor.Close(ctx)
	err := cursor.All(ctx, &departmentList)
	if err != nil {
		return nil, err
	}
	for _, v := range departmentList {
		cursor1, ctx1, _ := database.Db.QueryAll("department", bson.M{"pid": v.Id}, nil)
		err := cursor1.All(ctx1, &childrenDepartmentList)
		if err != nil {
			return nil, err
		}
		cursor1.Close(ctx)
		departmentList = append(departmentList, childrenDepartmentList...)
	}
	return departmentList, nil
}

//新建部门
func CreateDepartment(data interfaces.Department) error {
	err := database.Db.Insert("department", &data)
	if err != nil {
		return err
	}
	return nil
}

//删除部门
func RemoveDepartment(id string) error {
	err := database.Db.RemoveOne("department", bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

//移除人员去未分组
func RemoveHumanToUnclassified(ids []string) error {
	for _, id := range ids {
		err := database.Db.Update("human", bson.M{"pid": id}, bson.M{"$set": bson.M{"pid": "1", "dname": "未分组"}}, false)
		if err != nil {
			return err
		}
	}
	return nil
}

//注册人员
func RegisterHuman(human interfaces.Human) error {
	err := database.Db.Insert("human", &human)
	if err != nil {
		return err
	}
	err = database.Db.UpdateOne("countHuman", bson.M{}, bson.M{"$inc": bson.M{"count": 1}})
	if err != nil {
		return err
	}
	return nil
}

//批量删除人员
func RemoveHuman(ids []string) error {
	for _, id := range ids {
		err := database.Db.RemoveOne("human", bson.M{"_id": id})
		if err != nil {
			return err
		}
	}
	return nil
}

//编辑人员
func UpdateHuman(data interfaces.Human) error {
	if data.Path == "" {
		if err := database.Db.Update("human", bson.M{"_id": data.Id},
			bson.M{"$set": bson.M{"title": data.Title, "dname": data.Dname, "pid": data.Pid}}, false); err != nil {
			return err
		}
	} else {
		if err := database.Db.Update("human", bson.M{"_id": data.Id},
			bson.M{"$set": bson.M{"title": data.Title, "path": data.Path, "dname": data.Dname, "pid": data.Pid}}, false); err != nil {
			return err
		}
	}
	return nil
}

//获取用户ID
func GetHumanCount() (string, error) {
	res, err := database.Db.Query("countHuman", bson.M{})
	if err != nil {
		return "", err
	}
	id := res[0]["count"].(int32) + 1
	if id < 10 {
		return fmt.Sprintf("0%s", strconv.FormatInt(int64(id), 10)), nil
	} else if 10 <= id && id < 100 {
		return fmt.Sprintf("0%s", strconv.FormatInt(int64(id), 10)), nil
	} else {
		return fmt.Sprintf("%s", strconv.FormatInt(int64(id), 10)), nil
	}
}

//获取指定部门下人员信息
func QueryDepartmentHumanList(id string) ([]interfaces.Human, error) {
	var humanList []interfaces.Human
	opts := options.FindOptions{
		Sort: bson.M{"timestamp": -1},
	}
	cursor, ctx, _ := database.Db.QueryAll("human", bson.M{"pid": id}, &opts)
	defer cursor.Close(ctx)
	err := cursor.All(ctx, &humanList)
	if err != nil {
		return nil, err
	}
	return humanList, nil
}

//根据部门名称部门信息
func QuerySameDepartment(title string) ([]bson.M, error) {
	res, err := database.Db.Query("department", bson.M{"title": title})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//判断子部门下面是否还有字部门
func JudgeChildrenDepartment(id string) bool {
	var data interfaces.Department
	err := database.Db.QueryOne("department", bson.M{"_id": id}).Decode(&data)
	if err != nil {
		log.Error(err)
		return false
	}
	if data.Pid == "" {
		return true
	} else {
		return false
	}
}

//获取人员信息
func QueryHuman(id string, pid string, offset int64, limit int64) ([]interfaces.Human, int64, error) {
	var humanList []interfaces.Human
	var filter bson.M
	if id != "" {
		filter = bson.M{"_id": id}
	}
	if pid != "" {
		filter = bson.M{"pid": pid}
	}
	if id == "" && pid == "" {
		filter = bson.M{}
	}
	opts := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
		Sort:  bson.M{"timestamp": -1},
	}
	cursor, ctx, _ := database.Db.QueryAll("human", filter, &opts)
	defer cursor.Close(ctx)
	err := cursor.All(ctx, &humanList)
	if err != nil {
		return nil, 0, err
	}
	count, err := database.Db.QueryCount("human", filter)
	if err != nil {
		return nil, 0, err
	}
	return humanList, count, nil
}

//编辑部们
func UpdateDepartment(data interfaces.Department) error {
	err := database.Db.Update("department", bson.M{"_id": data.Id},
		bson.M{"$set": bson.M{"title": data.Title, "pid": data.Pid}}, false)
	if err != nil {
		return err
	}
	return nil
}

//根据用户列表获取用户
func QueryHumanByIds(data []string) ([]interfaces.Human, error) {
	var humanList []interfaces.Human
	var human interfaces.Human
	for i := 0; i < len(data); i++ {
		//Decode函数在*SingleResult为nil时, 查询结果为空直接抛异常
		err := database.Db.QueryOne("human", bson.M{"_id": data[i]}).Decode(&human)
		if err != nil {
			log.Error(err)
		}
		humanList = append(humanList, human)
	}
	return humanList, nil
}

//查询识别日志
func QueryIdentifyLog(humanType string, title string, starttime int64, endtime int64, offset int64, limit int64) ([]interfaces.Log, int64, error) {
	var data []interfaces.Log
	var filter bson.M
	if humanType == "" && title == "" {
		filter = bson.M{"timestamp": bson.M{"$gte": starttime, "$lte": endtime}}
	}
	if humanType != "" {
		filter = bson.M{"humanType": humanType, "timestamp": bson.M{"$gte": starttime, "$lte": endtime}}
	}
	if title != "" {
		filter = bson.M{"title": primitive.Regex{Pattern: title}, "timestamp": bson.M{"$gte": starttime, "$lte": endtime}}
	}
	opts := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
		Sort:  bson.M{"timestamp": -1},
	}
	cursor, ctx, _ := database.Db.QueryAll("log", filter, &opts)
	defer cursor.Close(ctx)
	err := cursor.All(ctx, &data)
	if err != nil {
		return nil, 0, err
	}
	count, err := database.Db.QueryCount("log", filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil
}
