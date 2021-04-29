package mq

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.jiangxingai.com/asp-hrm/cloud"
	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	"gitlab.jiangxingai.com/asp-hrm/router/mqtt/apiV1"
	log "k8s.io/klog"
)

type MqttClient struct {
	Broker   string
	Port     string
	Username string
	Password string
	ClientId string
}

func NewMqttClient() *MqttClient {
	return &MqttClient{
		Broker:   utils.GetEnv("MQTT_BROKER", "10.56.0.52"),
		Port:     utils.GetEnv("MQTT_PORT", "1883"),
		Username: utils.GetEnv("MQTT_USER", "hrm"),
		Password: utils.GetEnv("MQTT_PASSWORD", "123456"),
		ClientId: utils.GetEnv("MQTT_CLIENTID", ""),
	}
}

func connectCallback(client mqtt.Client) {
	log.Info("mqtt broker connect success")
	client.Subscribe("edge/+/latest", 2, apiV1.CaptureAnalysisHandler)
	client.Subscribe("hrm/face", 2, apiV1.GetAiResultHandler)
	client.Subscribe("+/rule/update", 2, apiV1.UpdateCameraHandler)
	client.Subscribe("asp/device/del", 2, apiV1.DeleteCameraHandler)
	log.Info("topic (re)subscribed")
}

func connectLostCallback(client mqtt.Client, err error) {
	log.Errorf("mqtt connect lost: %v", err)
	client.Unsubscribe("edge/+/latest", "hrm/face", "+/rule/update")
	ConnectMqtt()
}

func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func ConnectMqtt() {
	mqClient := NewMqttClient()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", mqClient.Broker, mqClient.Port))
	//opts.SetClientID(mqClient.ClientId)
	//opts.SetUsername(mqClient.Username)
	//opts.SetPassword(mqClient.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(60 * time.Second)
	opts.OnConnect = connectCallback
	opts.OnConnectionLost = connectLostCallback
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("mqtt broker connect fail: %s", token.Error())
		time.Sleep(time.Duration(1) * time.Second)
		ConnectMqtt()
	}
	cameraList := cloud.Setup()
	for _, camera := range cameraList {
		if camera.Enable == 1 {
			go apiV1.TaskCaptureRoutine(camera, client)
		}
	}
}
