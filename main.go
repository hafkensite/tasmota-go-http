package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var c mqtt.Client

type tasmotaConfig struct {
	IP              string `json:"ip"`
	MAC             string `json:"mac"`
	DescriptiveName string `json:"dn"`
	Topic           string `json:"t"`
	Software        string `json:"sw"`
	Relays          []int  `json:"rl"`
}

var configs map[string]tasmotaConfig
var statuses map[string]map[string]string

func main() {
	configs = make(map[string]tasmotaConfig)
	statuses = make(map[string]map[string]string)

	initTemplates()

	mqttBroker, found := os.LookupEnv("MQTT_BROKER")
	if !found {
		log.Fatal("MQTT_BROKER not set.")
	}
	mqttClientId, found := os.LookupEnv("MQTT_CLIENT_ID")
	if !found {
		log.Fatal("MQTT_CLIENT_ID not set.")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID(mqttClientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("stat/+/RESULT", 0, statHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("tasmota/discovery/+/config", 0, configHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	http.HandleFunc("/", rootHtml)
	http.HandleFunc("/manifest.json", webmanifestHandler)
	http.HandleFunc("/ws", websocketStatus)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
