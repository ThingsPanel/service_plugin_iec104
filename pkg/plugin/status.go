package plugin

import (
	"log/slog"
	"sync"
	"time"
)

type deviceStatus struct {
	DeviceId string
	ActiveCh chan struct{}
}

// Status 设备状态服务
type Status struct {
	upload    *uploader
	devices   map[string]*deviceStatus
	timeout   time.Duration // 离线超时时间
	mlock     sync.RWMutex  //数据锁
	lock      sync.Mutex    // 对象锁
	closed    bool          // 是否关闭
	closeCh   chan struct{}
	onlineCh  chan string
	offlineCh chan string
}

func newStatus(upload *uploader, timeout time.Duration) *Status {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &Status{
		upload:    upload,
		devices:   make(map[string]*deviceStatus),
		timeout:   timeout,
		closeCh:   make(chan struct{}),
		onlineCh:  make(chan string),
		offlineCh: make(chan string),
	}
}

// Active 激活设备状态
func (s *Status) Active(deviceId string) {
	device, loaded := s.load(deviceId)
	if !loaded {
		s.onlineCh <- deviceId
	}
	select {
	case device.ActiveCh <- struct{}{}:
	default:
	}
}

func (s *Status) Run() error {
	s.lock.Lock()
	if s.closed {
		s.lock.Unlock()
		return nil
	}
	s.lock.Unlock()

	go func() {
		ticker := time.NewTicker(s.timeout)
		defer ticker.Stop()
		for {
			select {
			case <-s.closeCh:
				return
			case <-ticker.C:
				s.mlock.RLock()
				ids := make([]string, 0, len(s.devices))
				for _, item := range s.devices {
					select {
					case <-item.ActiveCh:
					default:
						ids = append(ids, item.DeviceId)
					}
				}
				s.mlock.RUnlock()

				if len(ids) > 0 {
					s.mlock.Lock()
					for _, id := range ids {
						delete(s.devices, id)
						s.offlineCh <- id
					}
					s.mlock.Unlock()
				}
			}
		}
	}()

	for {
		select {
		case <-s.closeCh:
			return nil
		case id := <-s.onlineCh:
			err := s.upload.UploadStatus(id, true, true)
			if err != nil {
				slog.Warn("upload status", "error", err)
			}
		case id := <-s.offlineCh:
			err := s.upload.UploadStatus(id, false, true)
			if err != nil {
				slog.Warn("upload status", "error", err)
			}
		}
	}
}

func (s *Status) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.closed {
		return nil
	}

	close(s.closeCh)
	s.closed = true

	s.mlock.Lock()
	defer s.mlock.Unlock()

	for _, item := range s.devices {
		_ = s.upload.UploadStatus(item.DeviceId, false)
	}

	return nil
}

func (s *Status) load(deviceId string) (device *deviceStatus, isExist bool) {
	s.mlock.RLock()
	val, ok := s.devices[deviceId]
	if ok {
		s.mlock.RUnlock()
		return val, true
	}
	s.mlock.RUnlock()

	s.mlock.Lock()
	defer s.mlock.Unlock()

	val = &deviceStatus{
		DeviceId: deviceId,
		ActiveCh: make(chan struct{}),
	}
	s.devices[deviceId] = val
	return val, false
}
