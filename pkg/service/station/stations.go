package station

import "sync"

type Stations struct {
	list map[uint16]*Station
	lock sync.Mutex
}

func NewStations() *Stations {
	return &Stations{
		list: make(map[uint16]*Station),
	}
}

// Load 载入站信息
func (s *Stations) Load(addr uint16) *Station {
	s.lock.Lock()
	defer s.lock.Unlock()

	row, ok := s.list[addr]
	if ok {
		return row
	}

	row = NewStation(addr)
	s.list[addr] = row
	return row
}

// Exist 是否存在
func (s *Stations) Exist(addr uint16) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.list[addr]
	return ok
}

func (s *Stations) List() []*Station {
	s.lock.Lock()
	defer s.lock.Unlock()

	var res []*Station
	for _, item := range s.list {
		res = append(res, item)
	}
	return res
}

func (s *Stations) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.list)
}
