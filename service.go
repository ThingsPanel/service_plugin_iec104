package iec104_master

import (
	"context"
	"iec104-slave/pkg/plugin"
	"iec104-slave/pkg/plugin/tpapi"
	"iec104-slave/pkg/plugin/voucher"
	"iec104-slave/pkg/service"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	services map[string]*service.Service
	plugin   *plugin.Server
	closeCh  chan struct{}
	closed   bool
	dlock    sync.Mutex // 数据锁
	rlock    sync.Mutex // 状态锁
}

func NewService(plugin *plugin.Server) *Service {
	return &Service{
		services: make(map[string]*service.Service),
		plugin:   plugin,
		closeCh:  make(chan struct{}),
	}
}

// Init 初始化服务列表,从远程获取服务列表
func (s *Service) Init() error {
	list, err := s.plugin.Api().GetServiceList()
	if err != nil {
		return err
	}

	// 去除不存在还是跑的服务
	s.dlock.Lock()
	for _, item := range s.services {
		var found bool
		for _, target := range list {
			vo, err := target.Vouchers()
			if err != nil {
				continue
			}
			if item.Remote() == vo.Remote() {
				found = true
			}
		}
		if !found {
			if err := item.Close(); err != nil {
			}
			delete(s.services, item.Remote())
		}
	}
	s.dlock.Unlock()

	for _, item := range list {
		vo, err := item.Vouchers()
		if err != nil {
			slog.Error("parse voucher", "data", item.Voucher, "err", err)
		} else {
			s.Load(vo)
		}
	}
	return nil
}

// Load 获取服务不存在则创建
func (s *Service) Load(vo *voucher.Service) *service.Service {
	s.dlock.Lock()
	defer s.dlock.Unlock()

	svc, ok := s.services[vo.Remote()]
	if ok {
		return svc
	}

	svc = service.NewService(service.Config{
		Remote:  vo.Remote(),
		Handler: s.onReceived,
	})
	s.services[vo.Remote()] = svc

	return svc
}

func (s *Service) Exist(vo *voucher.Service) bool {
	s.dlock.Lock()
	defer s.dlock.Unlock()
	_, ok := s.services[vo.Remote()]
	return ok
}

// Restart 重新加载服务数据
func (s *Service) Restart() {
	s.rlock.Lock()
	defer s.rlock.Unlock()
	err := s.Init()
	if err != nil {
		slog.Warn("restart service", "err", err)
	}
}

// Run 运行服务列表,定时检查
func (s *Service) Run() error {
	// 第一次运行获取服务列表
	if err := s.Init(); err != nil {
		return err
	}

	var runSvc = func() {
		s.dlock.Lock()
		defer s.dlock.Unlock()
		for _, item := range s.services {
			if !item.IsRunning() {
				go func() {
					if err := item.Run(); err != nil {
						slog.Error("service", "addr", item.Remote(), "err", err)
					}
				}()
			}
		}
	}

	runSvc()

	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-s.closeCh:
			return nil
		case <-ticker.C:
			runSvc()
		}
	}
}

func (s *Service) Close() error {
	s.rlock.Lock()
	defer s.rlock.Unlock()
	close(s.closeCh)
	return nil
}

func (s *Service) onReceived(ctx context.Context, station uint16, ioa uint, data any) error {
	slog.Info("onReceived", "station", station, "ioa", ioa, "data", data)

	dev, err := s.getDeviceBySid(station)
	if err != nil {
		return err
	}

	s.plugin.Active(dev.ID)

	attrs := make(map[string]any)
	attrs[strconv.Itoa(int(ioa))] = data

	return s.plugin.Upload().PostAttrs(dev.ID, attrs)
}

func (s *Service) getDeviceBySid(sid uint16) (*tpapi.DeviceConfig, error) {
	return s.plugin.Api().GetDeviceByNumber(strconv.Itoa(int(sid)))
}
