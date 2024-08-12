package plugin

import (
	"testing"
	"watersql/pkg/config"
)

func Test_Mqtt(t *testing.T) {
	mqtt := NewMqttClient(config.Mqtt{
		Broker:   "47.92.253.145:1883",
		Username: "root",
		Password: "root",
	})
	err := mqtt.Connect()
	if err != nil {
		t.Fatal(err)
	}
}
