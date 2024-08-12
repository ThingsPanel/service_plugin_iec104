package station

import "sync"

type Devices struct {
	list map[uint]*Device
	lock sync.Mutex
}

func NewDevices() *Devices {
	return &Devices{
		list: make(map[uint]*Device),
	}
}

func (s *Devices) Set(d *Device) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list[d.Address] = d
}

func (s *Devices) Add(addr uint) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list[addr] = NewDevice(addr)
}

func (s *Devices) Get(addr uint) (*Device, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	row, ok := s.list[addr]
	return row, ok
}

func (s *Devices) Del(addr uint) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.list, addr)
}

func (s *Devices) Load(addr uint) *Device {
	s.lock.Lock()
	defer s.lock.Unlock()
	row, ok := s.list[addr]
	if ok {
		return row
	}
	row = NewDevice(addr)
	s.list[addr] = row
	return row
}

func (s *Devices) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.list)
}
