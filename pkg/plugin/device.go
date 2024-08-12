package plugin

// Device 设备信息
type Device interface {
	DeviceNumber() string
	DeviceName() string
	Description() string
}
