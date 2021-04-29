package e

var MsgFlags = map[int]string{
  SUCCESS:        "success",
  ERROR:          "fail",
  INVALID_PARAMS: "请求参数错误",

  ERROR_DEPARTMENT_NAME: "部门已存在",
  ERROR_SAVE_IMAGE:      "保存图片失败",
  ERROR_IMAGE_NAME:      "请上传正确图片格式",
  ERROR_REGISTER_HUMAN:  "百度sdk注册人员失败",
  ERROR_UPDATE_HUMAN:    "百度sdk更新人员失败",
  ERROR_REMOVE_HUMAN:    "百度sdk删除人员失败",
  ERROR_MUTLI_HUMAN:     "上传失败,检测到图片中存在多张人脸信息",
  ERROR_QUALITY_HUMAN:   "上传失败,图片中未检测到人脸信息",
  ERROR_EXTRA_HUMAN:     "上传失败,请重试",
  ERROR_SAME_HUMAN:      "此人已在系统中,请勿重复录入",
}

func GetMsg(code int) string {
  msg, ok := MsgFlags[code]
  if ok {
    return msg
  }
  return MsgFlags[ERROR]
}
