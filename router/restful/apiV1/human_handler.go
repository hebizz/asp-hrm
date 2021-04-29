package apiV1

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	e "gitlab.jiangxingai.com/asp-hrm/pkg/response"
	"gitlab.jiangxingai.com/asp-hrm/pkg/sdk"
	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "k8s.io/klog"
)

//获取部门树
func GetDepartmentHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var fatherDepartment, childrenDepartment, grandsonDepartment []interfaces.DepartmentTree
	var lastDepartment interfaces.DepartmentTree
	res, err := QueryDepartmentList([]interfaces.DepartmentTree{})
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, nil, nil)
		return
	}
	for _, v := range res.([]interfaces.DepartmentTree) {
		if v.Pid == "" {
			//未分组放到最后
			if v.Id != "1" {
				fatherDepartment = append(fatherDepartment, v)
			} else {
				lastDepartment = v
			}
		} else if JudgeChildrenDepartment(v.Pid) {
			childrenDepartment = append(childrenDepartment, v)
		} else {
			grandsonDepartment = append(grandsonDepartment, v)
		}
	}
	fatherDepartment = append(fatherDepartment, lastDepartment)
	for l, k := range childrenDepartment {
		for x, y := range grandsonDepartment {
			if y.Pid == k.Id {
				grandsonDepartment[x].HList, _ = QueryDepartmentHumanList(y.Id)
				childrenDepartment[l].Children = append(childrenDepartment[l].Children, grandsonDepartment[x])
			}
		}
	}
	for p, v := range fatherDepartment {
		fatherDepartment[p].HList, _ = QueryDepartmentHumanList(v.Id)
		for l, k := range childrenDepartment {
			if k.Pid == v.Id {
				childrenDepartment[l].HList, _ = QueryDepartmentHumanList(k.Id)
				fatherDepartment[p].Children = append(fatherDepartment[p].Children, childrenDepartment[l])
			}
		}
	}
	log.Info(fatherDepartment)
	app.Response(http.StatusOK, e.SUCCESS, nil, fatherDepartment)
}

//获取所有部门列表
func GetDepartmentListHandler(c *gin.Context) {
	app := e.Gin{C: c}
	res, err := QueryDepartmentList([]interfaces.Department{})
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, nil, nil)
		return
	}
	log.Info(res)
	app.Response(http.StatusOK, e.SUCCESS, nil, res)
}

//新建部门
func CreateDepartmentHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.Department
	err := c.BindJSON(&data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	res, err := QuerySameDepartment(data.Title)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	if len(res) > 0 {
		app.Response(http.StatusBadRequest, e.ERROR_DEPARTMENT_NAME, nil, nil)
		return
	}
	data.Timestamp = time.Now().Unix()
	data.Id = primitive.NewObjectID().Hex()
	data.Type = "department"
	err = CreateDepartment(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//编辑部们
func UpdateDepartmentHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.Department
	err := c.BindJSON(&data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	err = UpdateDepartment(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//删除部门
func DeleteDepartmentHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.Department
	var departmentIdList []string
	err := c.BindJSON(&data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	departmentIdList = append(departmentIdList, data.Id)
	err = RemoveDepartment(data.Id)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	departmentMsg, err := GetChildrenDepartment(data.Id)
	for _, v := range departmentMsg {
		departmentIdList = append(departmentIdList, v.Id)
		err = RemoveDepartment(v.Id)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
			return
		}
	}
	err = RemoveHumanToUnclassified(departmentIdList)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//获取人员信息
func GetHumanHandler(c *gin.Context) {
	app := e.Gin{C: c}
	id := c.Query("id")
	pid := c.Query("pid")
	offset, err1 := strconv.ParseInt(c.Query("offset"), 10, 64)
	limit, err2 := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err1 != nil || err2 != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	humanList, count, err := QueryHuman(id, pid, offset, limit)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	log.Info(humanList)
	app.Response(http.StatusOK, e.SUCCESS, nil, gin.H{"data": humanList, "count": count})
}

//人员注册
func RegisterHumanHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.Human
	data.Title = c.PostForm("title")
	data.Pid = c.PostForm("pid")
	data.Dname = c.PostForm("dname")
	image, err := c.FormFile("image")
	if err != nil || data.Title == "" || data.Pid == "" || data.Dname == "" {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	imgTmp := strings.Split(image.Filename, ".")
	supportImg := []string{"jpeg", "jpg", "png"}
	if len(imgTmp) < 2 || !utils.IsExistItem(imgTmp[1], supportImg) {
		app.Response(http.StatusBadRequest, e.ERROR_IMAGE_NAME, nil, nil)
		return
	}
	data.Timestamp = time.Now().Unix()
	timestamp := strconv.FormatInt(data.Timestamp, 10)
	data.Path = fmt.Sprintf("%s%s.jpeg", FrontendPath, timestamp)
	data.Rpath = fmt.Sprintf("%s%s.jpeg", ImageSavePath, timestamp)
	if !utils.Exists(ImageSavePath) {
		if err := os.MkdirAll(ImageSavePath, os.ModePerm); err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
			return
		}
	}
	err = c.SaveUploadedFile(image, data.Rpath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	flag, err := sdk.CheckMultiHumanOnBaiduSdk(data.Rpath)
	if err != nil {
		app.Response(http.StatusBadRequest, e.ERROR_EXTRA_HUMAN, err, nil)
		return
	}
	if flag != 0 {
		app.Response(http.StatusBadRequest, flag, err, nil)
		return
	}

	flag, err = sdk.CheckHumanOnBaiduSdk(data.Rpath)
	if err != nil {
		app.Response(http.StatusBadRequest, e.ERROR_EXTRA_HUMAN, err, nil)
		return
	}
	if flag != 0 {
		app.Response(http.StatusBadRequest, flag, err, nil)
		return
	}

	humanData, err := GenerateHumanData(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	humanData.FaceToken, err = sdk.RegisterHumanToBaiduSdk("register", humanData)
	log.Info(data.FaceToken)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_REGISTER_HUMAN, err, nil)
		return
	}
	err = RegisterHuman(humanData)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//人员删除
func RemoveHumanHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.HumanIds
	err := c.BindJSON(&data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	humanList, err := QueryHumanByIds(data.Ids)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	log.Info(humanList)
	for _, human := range humanList {
		err = sdk.RemoveHumanToBaiduSdk(human)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR_REMOVE_HUMAN, err, nil)
			return
		}
	}
	err = RemoveHuman(data.Ids)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//编辑人员
func UpdateHumanHandler(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.Human
	data.Id = c.PostForm("id")
	data.Title = c.PostForm("title")
	data.Pid = c.PostForm("pid")
	data.Dname = c.PostForm("dname")
	image, err := c.FormFile("image")
	if data.Id == "" || data.Title == "" || data.Dname == "" || data.Pid == "" {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	if err != nil {
		data.Path = ""
	} else {
		imgTmp := strings.Split(image.Filename, ".")
		supportImg := []string{"jpeg", "jpg", "png"}
		if len(imgTmp) < 2 || !utils.IsExistItem(imgTmp[1], supportImg) {
			app.Response(http.StatusBadRequest, e.ERROR_IMAGE_NAME, nil, nil)
			return
		}
		data.Timestamp = time.Now().Unix()
		data.Rpath = fmt.Sprintf("%s%s.%s", ImageSavePath, strconv.FormatInt(data.Timestamp, 10), imgTmp[1])
		data.Path = fmt.Sprintf("%s%s.%s", FrontendPath, strconv.FormatInt(data.Timestamp, 10), imgTmp[1])
		err = c.SaveUploadedFile(image, data.Rpath)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR_SAVE_IMAGE, err, nil)
			return
		}
		flag, err := sdk.CheckMultiHumanOnBaiduSdk(data.Rpath)
		if err != nil {
			app.Response(http.StatusBadRequest, e.ERROR_EXTRA_HUMAN, err, nil)
			return
		}
		if flag != 0 {
			app.Response(http.StatusBadRequest, flag, err, nil)
			return
		}
		data.FaceToken, err = sdk.RegisterHumanToBaiduSdk("update", data)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR_EXTRA_HUMAN, err, nil)
			return
		}
	}
	err = UpdateHuman(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//获取识别日志
func GetIdentifyLogHandler(c *gin.Context) {
	app := e.Gin{C: c}
	humanType := c.Query("humanType")
	title := c.Query("title")
	starttime, err1 := strconv.ParseInt(c.Query("starttime"), 10, 64)
	endtime, err2 := strconv.ParseInt(c.Query("endtime"), 10, 64)
	offset, err3 := strconv.ParseInt(c.Query("offset"), 10, 64)
	limit, err4 := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	res, count, err := QueryIdentifyLog(humanType, title, starttime, endtime, offset, limit)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, gin.H{"data": res, "total": count})
}
