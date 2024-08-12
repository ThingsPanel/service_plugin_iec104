package iec104_master

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"iec104-slave/pkg/config"
	"iec104-slave/pkg/plugin"
	"iec104-slave/pkg/plugin/voucher"
	"log/slog"
)

type Server struct {
	plugin  *plugin.Server
	config  config.Config
	service *Service
}

func NewServer(config config.Config) *Server {
	svr := &Server{
		config: config,
	}
	svr.plugin = plugin.NewServer(config.Things, svr)
	svr.service = NewService(svr.plugin)
	return svr
}

// OnRequestDevices 获取设备列表
func (s *Server) OnRequestDevices(ctx context.Context, voucher *voucher.Service, page int, size int) ([]plugin.Device, int, error) {
	slog.Debug("devices request", "voucher", voucher, "page", page, "size", size)

	svc := s.service.Load(voucher)
	list := svc.Stations()

	result := make([]plugin.Device, 0, len(list))
	for _, item := range list {
		fmt.Println("station", item)
	}

	start := page - 1
	end := page * size
	if end > len(result) {
		return result, len(result), nil
	}
	if start > len(result) {
		return []plugin.Device{}, 0, nil
	}
	return result[start:end], len(result), nil
}

// OnEventNotify 事件通知
func (s *Server) OnEventNotify(ctx context.Context) {
	s.service.Restart()
}

// OnDisconnectDevice 断开设备链接
func (s *Server) OnDisconnectDevice(ctx context.Context, deviceId string) error {
	return nil
}

func (s *Server) Run() error {
	var eg errgroup.Group
	eg.Go(func() error {
		return s.plugin.Run()
	})
	eg.Go(func() error {
		return s.service.Run()
	})
	return eg.Wait()
}
