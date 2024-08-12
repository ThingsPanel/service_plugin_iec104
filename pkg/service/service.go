package service

import (
	"context"
	"errors"
	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
	"iec104-slave/pkg/service/station"
	"log/slog"
	"sync"
	"time"
)

type Service struct {
	client  *cs104.Client
	station *station.Stations
	config  Config
	closed  bool
	running bool
	closeCh chan struct{}
	lock    sync.Mutex
	ctx     context.Context
}

func NewService(c Config) *Service {
	svc := &Service{
		config:  c,
		station: station.NewStations(),
		ctx:     c.Context,
	}

	if svc.ctx == nil {
		svc.ctx = context.Background()
	}
	return svc
}

// Stations 设备列表
func (s *Service) Stations() []*station.Station {
	return s.station.List()
}

// Remote 从站地址
func (s *Service) Remote() string {
	return s.config.Remote
}

func (s *Service) ServiceId() string {
	return s.config.ServiceId
}

func (s *Service) Init() error {
	option := cs104.NewOption()
	option.SetConfig(cs104.Config{
		ConnectTimeout0:   5 * time.Second,
		SendUnAckLimitK:   8,
		SendUnAckTimeout1: 15 * time.Second,
		RecvUnAckLimitW:   1,
		RecvUnAckTimeout2: 10 * time.Second,
		IdleTimeout3:      20 * time.Second,
	})

	if err := option.AddRemoteServer(s.config.Remote); err != nil {
		return err
	}

	client := cs104.NewClient(newIecClient(s.station, s.onDataReceive), option)
	client.LogMode(false)
	client.SetOnConnectHandler(s.onConnectHandler)
	client.SetConnectionLostHandler(s.onConnectionLostHandler)
	client.SetActiveHandler(s.onActiveHandler)

	s.client = client

	return nil
}

func (s *Service) Run() error {
	if s.IsRunning() {
		return errors.New("service is running")
	}

	// 初始化客户端
	if err := s.Init(); err != nil {
		return err
	}

	// 初始化状态
	s.closeCh = make(chan struct{})
	s.closed = false
	s.running = true

	// 开启tcp客户端
	err := s.client.Start()
	if err != nil {
		return err
	}

	// 检查是否链接成功
	time.Sleep(time.Second * 5)
	if !s.client.IsConnected() {
		return errors.New("client connect failed")
	}

	ticker := time.NewTicker(time.Minute)
	defer func() {
		s.lock.Lock()
		s.closed = true
		s.running = false
		s.lock.Unlock()

		ticker.Stop()
		s.client.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-s.closeCh:
			return nil
		case <-ticker.C:
			if !s.client.IsConnected() {
				return errors.New("client closed")
			}
			if err := s.sendInterrogation(1); err != nil {
				return err
			}
		}
	}
}

func (s *Service) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	close(s.closeCh)
	return s.client.Close()
}

func (s *Service) IsRunning() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.running
}

// 发送startDt激活指令
func (s *Service) onConnectHandler(c *cs104.Client) {
	slog.Warn("cs104 client has connected")
	go func() {
		c.SendStartDt()
	}()
}

func (s *Service) onConnectionLostHandler(c *cs104.Client) {
	slog.Warn("cs104 client lost connection")
	s.Close()
}

// 激活后发送总召唤
func (s *Service) onActiveHandler(c *cs104.Client) {
	go func() {
		s.sendInterrogation(1)
	}()
}

// 遥测数据回调
func (s *Service) onDataReceive(station uint16, ioa uint, data any) {
	s.config.Handler(context.TODO(), station, ioa, data)
}

// 发送总召唤
func (s *Service) sendInterrogation(addr uint16) error {
	err := s.client.InterrogationCmd(
		asdu.CauseOfTransmission{
			Cause: asdu.Activation,
		},
		asdu.CommonAddr(addr),
		asdu.QOIStation,
	)
	if err != nil {
		slog.Warn("send interrogation", "err", err)
	}
	return err
}
