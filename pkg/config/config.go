package config

import (
  log "k8s.io/klog"
)

var (
  BuildVersion string
  BuildTime    string
  BuildHash    string
  GoVersion    string
  GoBuildType  string
)


const (
  HttpPort = ":9006"
)

//打印版本信息
func Setup() {
  log.Infof("Version: %s\nBuild Time: %s\nBuild Hash: %s\nGo Version: %s\nGoBuildType: %s\n",
    BuildVersion, BuildTime, BuildHash, GoVersion, GoBuildType)
}



