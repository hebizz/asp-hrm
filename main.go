package main

import (
  "github.com/gin-gonic/gin"
  "gitlab.jiangxingai.com/asp-hrm/database"
  "gitlab.jiangxingai.com/asp-hrm/pkg/config"
  "gitlab.jiangxingai.com/asp-hrm/router/mqtt"
  "gitlab.jiangxingai.com/asp-hrm/router/restful"
  log "k8s.io/klog"
)

func init() {
  config.Setup()
  database.Setup()
  go mq.ConnectMqtt()
}

func main() {
  gin.SetMode("debug")
  engine := gin.Default()
  engine.MaxMultipartMemory = 10
  restful.InitRestFul(engine)
  if err := engine.Run(config.HttpPort); err != nil {
    log.Info(err)
  }
}
