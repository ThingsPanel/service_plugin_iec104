package plugin

import (
	"encoding/json"
	"errors"
	"log/slog"
)

type uploader struct {
	mqtt    *MqttClient
	attrCh  chan []byte
	statCh  chan statusData
	closeCh chan struct{}
}

func newUploader(mqtt *MqttClient) *uploader {
	return &uploader{
		mqtt:    mqtt,
		attrCh:  make(chan []byte, 10),
		statCh:  make(chan statusData, 10),
		closeCh: make(chan struct{}),
	}
}

// PostAttrs 异步上传，放到发送队列
func (s *uploader) PostAttrs(deviceId string, attrs map[string]any) error {
	payload, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	obj := map[string]any{
		"device_id": deviceId,
		"values":    payload,
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	select {
	case s.attrCh <- data:
	default:
	}
	return errors.New("post attr chan fulled")
}

// UploadAttrs 同步上传
func (s *uploader) UploadAttrs(deviceId string, attrs map[string]any) error {
	payload, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	obj := map[string]any{
		"device_id": deviceId,
		"values":    payload,
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	token := s.mqtt.Publish("devices/telemetry", data)
	if token.Wait() && token.Error() != nil {
		slog.Warn("upload attrs", "error", token.Error())
		return token.Error()
	}

	return nil
}

// PostStatus 异步上传状态
func (s *uploader) PostStatus(deviceId string, online bool) error {
	data := statusData{
		deviceId: deviceId,
		online:   online,
	}
	select {
	case s.statCh <- data:
	default:
	}
	return errors.New("status chan fulled")
}

// UploadStatus 同步上传状态
func (s *uploader) UploadStatus(deviceId string, online bool, async ...bool) error {
	status := "1"
	if !online {
		status = "0"
	}

	token := s.mqtt.Publish("devices/status/"+deviceId, []byte(status))
	if len(async) > 0 && async[0] {
		return token.Error()
	} else {
		if token.Wait() && token.Error() != nil {
			slog.Warn("upload status", "error", token.Error())
			return token.Error()
		}
	}
	return nil
}

func (s *uploader) Run() error {
	for {
		select {
		case <-s.closeCh:
			return nil
		case data := <-s.attrCh:
			s.mqtt.Publish("devices/telemetry", data)
		case status := <-s.statCh:
			s.UploadStatus(status.deviceId, status.online)
		}
	}
}

func (s *uploader) Close() error {
	close(s.closeCh)
	return nil
}

type statusData struct {
	deviceId string
	online   bool
}
