package utils

import (
  "bytes"
  "crypto/md5"
  "encoding/base64"
  "encoding/hex"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "os/exec"
  "reflect"
  "strconv"
  "strings"

  log "k8s.io/klog"
)


func GetEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}


func CommandExecuteLogs(cmd *exec.Cmd) (string, error) {
  var out bytes.Buffer
  var stderr bytes.Buffer
  cmd.Stdout = &out
  cmd.Stderr = &stderr
  err := cmd.Run()
  if err != nil {
    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    return stderr.String(), errors.New(stderr.String())
  }
  log.Info("cmd execute result is :\n" + out.String())
  return out.String(), nil
}

func CommandExecuteWithAllLogs(command string) (string, error) {
  cmd := exec.Command("bash", "-c", command)
  output, err := cmd.CombinedOutput()
  log.Info("cmd execute result is :\n" + string(output))
  return string(output), err
}

func GenerateMd5(rawStr string) string {
  data := []byte(rawStr)
  md5Str := fmt.Sprintf("%x", md5.Sum(data))
  return md5Str
}

func GetMD5Str(timestamp int64) string {
  str := strconv.FormatInt(timestamp, 10)
  h := md5.New()
  h.Write([]byte(str))
  return hex.EncodeToString(h.Sum(nil))
}

func IsExistItem(value interface{}, array interface{}) bool {
  switch reflect.TypeOf(array).Kind() {
  case reflect.Slice:
    s := reflect.ValueOf(array)
    for i := 0; i < s.Len(); i++ {
      if reflect.DeepEqual(value, s.Index(i).Interface()) {
        return true
      }
    }
  }
  return false
}

func Exists(path string) bool {
  _, err := os.Stat(path) //os.Stat获取文件信息
  if err != nil {
    if os.IsExist(err) {
      return true
    }
    return false
  }
  return true
}

func HttpGet(url string) (string, error) {
  resp, err := http.Get(url)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  return string(body), nil
}

func HttpPost(url string, formData string) (string, error) {
  resp, err := http.Post(url, "application/x-www-form-urlencoded",
    strings.NewReader(formData))
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  return string(body), nil
}

func HttpPostJson(url string, jsonData []byte) (string, error) {
  resp, err := http.Post(url, "application/json",
    bytes.NewBuffer(jsonData))
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  return string(body), nil
}


func Base64ToFile(filename int64, base64Str string) (string, error) {
  err := os.MkdirAll("/data/local/asp/media/face/", os.ModePerm)
  if err != nil {
    return "", err
  }
  str, err := base64.StdEncoding.DecodeString(base64Str)
  if err != nil {
    return "", err
  }
  path := fmt.Sprintf("/data/local/asp/media/face/%d.jpg", filename)
  err = ioutil.WriteFile(path, str, 0666)
  if err != nil {
    return "", err
  }
  return path, nil
}