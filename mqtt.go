package main

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var configHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var t tasmotaConfig
	err := json.Unmarshal(msg.Payload(), &t)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed :%v\n", t)

	configs[t.Topic] = t

	for idx, v := range t.Relays {
		if v == 1 {
			statusTopic := fmt.Sprintf("cmnd/%s/POWER%d", t.Topic, idx+1)
			fmt.Printf("Reqesting status %s\n", statusTopic)
			if token := c.Publish(statusTopic, 0, false, ""); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}
	}
}

var statHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	t := strings.Split(msg.Topic(), "/")[1]
	fmt.Printf("TOPIC: %s\n", t)

	statusMap := make(map[string]string)
	err := json.Unmarshal(msg.Payload(), &statusMap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed :%v\n", statusMap)
	deviceStatus, ok := statuses[t]
	if !ok {
		deviceStatus = make(map[string]string)
		statuses[t] = deviceStatus
	}
	for k, v := range statusMap {
		deviceStatus[k] = v
	}
	for conn := range activeConns {
		conn.WriteJSON(statuses)
	}
	fmt.Printf("statuses :%v\n", statuses)
}
