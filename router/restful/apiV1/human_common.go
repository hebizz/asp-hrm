package apiV1

import (
  "strings"

  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/Chain-Zhang/pinyin"
  "gitlab.jiangxingai.com/asp-hrm/interfaces"
)

const (
  ImageSavePath = "/data/local/asp/media/face/"
  FrontendPath = "/api/media/face/"
)

//生成人员信息
func GenerateHumanData(data interfaces.Human) (interfaces.Human, error) {
  data.Id = primitive.NewObjectID().Hex()
  data.Type = "human"
  str, err := pinyin.New(data.Title).Split(" ").Mode(pinyin.WithoutTone).Convert()
  if err != nil {
    return interfaces.Human{}, err
  }
  data.UserId = "ID"
  strList := strings.Split(str, " ")
  for i := 0; i < len(strList); i++ {
    data.UserId = data.UserId + string(strList[i][0])
  }
  count, err := GetHumanCount()
  if err != nil {
    return interfaces.Human{}, err
  }
  data.UserId = data.UserId + count
  return data, nil
}
