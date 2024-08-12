package server

type Device struct {
	DeviceNumber string `json:"device_number"`
	DeviceName   string `json:"device_name"`
	Description  string `json:"description"`
}
