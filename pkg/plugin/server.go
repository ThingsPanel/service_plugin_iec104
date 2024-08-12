package plugin

import (
	_ "embed"
	"fmt"
	"golang.org/x/sync/errgroup"
	"iec104-slave/pkg/config"
	"iec104-slave/pkg/plugin/server"
	"iec104-slave/pkg/plugin/tpapi"
	"log/slog"
	"time"
)

type Server struct {
	http    *server.Server
	upload  *uploader
	mqtt    *MqttClient
	status  *Status
	tpapi   *tpapi.Client
	config  config.Things
	closeCh chan struct{}
}

func NewServer(config config.Things, provider Provider) *Server {
	srv := &Server{
		config:  config,
		closeCh: make(chan struct{}),
	}
	srv.mqtt = NewMqttClient(config.Mqtt)
	srv.http = server.NewServer(newHandler(provider))
	srv.upload = newUploader(srv.mqtt)
	srv.status = newStatus(srv.upload, time.Duration(config.Device.Timeout)*time.Minute)
	srv.tpapi = tpapi.NewClient(config.Remote, "IEC104")
	return srv
}

func (s *Server) Upload() *uploader {
	return s.upload
}

// Active 激活设备状态
func (s *Server) Active(deviceId string) {
	s.status.Active(deviceId)
}

// Api api客户端
func (s *Server) Api() *tpapi.Client {
	return s.tpapi
}

func (s *Server) Run() error {
	if err := s.mqtt.Connect(); err != nil {
		return fmt.Errorf("mqtt connect fail: %w", err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		return s.upload.Run()
	})
	eg.Go(func() error {
		return s.status.Run()
	})
	eg.Go(func() error {
		return s.http.Run(s.config.Local)
	})
	eg.Go(func() error {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		_ = s.tpapi.SendHeartbeat()

		for {
			select {
			case <-s.closeCh:
				return nil
			case <-ticker.C:
				_ = s.tpapi.SendHeartbeat()
			}
		}
	})

	slog.Info("http server run", "addr", s.config.Local)

	return eg.Wait()
}

func (s *Server) Close() error {
	close(s.closeCh)

	if err := s.status.Close(); err != nil {
		return err
	}
	if err := s.upload.Close(); err != nil {
		return err
	}
	if err := s.mqtt.Disconnect(); err != nil {
		return err
	}
	if err := s.http.Close(); err != nil {
		return err
	}
	return nil
}
