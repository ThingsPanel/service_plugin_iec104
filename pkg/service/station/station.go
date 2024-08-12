package station

import (
	"fmt"
	"strconv"
)

type Station struct {
	Address uint16   // 站地址
	Devices *Devices // 子设备
}

func NewStation(address uint16) *Station {
	return &Station{
		Address: address,
		Devices: NewDevices(),
	}
}

func (s Station) DeviceNumber() string {
	return strconv.Itoa(int(s.Address))
}

func (s Station) DeviceName() string {
	return fmt.Sprintf("子设备数量%d", s.Devices.Count())
}

func (s Station) Description() string {
	return fmt.Sprintf("子设备数量%d", s.Devices.Count())
}
