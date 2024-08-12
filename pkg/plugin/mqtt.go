package plugin

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"iec104-slave/pkg/config"
	"log/slog"
	"time"
)

type MqttClient struct {
	client   mqtt.Client
	address  string
	username string
	password string
}

func NewMqttClient(config config.Mqtt) *MqttClient {
	mq := &MqttClient{
		address:  config.Broker,
		username: config.Username,
		password: config.Password,
	}

	opts := mqtt.NewClientOptions()
	opts.SetClientID(uuid.New().String())
	opts.AddBroker(config.Broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(time.Second * 3)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		slog.Info("mqtt", "status", "connected")
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		slog.Warn("mqtt", "status", "reconnect")
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		slog.Warn("mqtt", "status", "disconnect")
	})

	mq.client = mqtt.NewClient(opts)

	return mq
}

func (s *MqttClient) Connect() error {
	if s.client.IsConnected() {
		return nil
	}
	token := s.client.Connect()
	token.Wait()
	return token.Error()
}

func (s *MqttClient) IsConnect() bool {
	if s.client == nil {
		return false
	}
	return s.client.IsConnected()
}

func (s *MqttClient) Disconnect() error {
	if s.client == nil {
		return nil
	}
	s.client.Disconnect(0)
	return nil
}

func (s *MqttClient) Publish(topic string, payload []byte, qos ...byte) mqtt.Token {
	var qs byte
	if len(qos) > 0 {
		qs = qos[0]
	}
	return s.client.Publish(topic, qs, false, payload)
}

func (s *MqttClient) Subscribe(topic string, callback func(mqtt.Message), qos ...byte) mqtt.Token {
	var qs byte
	if len(qos) > 0 {
		qs = qos[0]
	}
	return s.client.Subscribe(topic, qs, func(client mqtt.Client, message mqtt.Message) {
		if callback != nil {
			callback(message)
		}
	})
}
