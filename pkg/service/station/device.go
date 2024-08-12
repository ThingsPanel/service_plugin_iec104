package station

type Device struct {
	Address uint
	Type    DeviceType
}

func NewDevice(addr uint) *Device {
	return &Device{
		Address: addr,
	}
}
