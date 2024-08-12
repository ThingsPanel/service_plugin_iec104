package tpapi

import "time"

// Device thingsPanel设备信息
type Device struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Voucher         string      `json:"voucher"`
	TenantID        string      `json:"tenant_id"`
	IsEnabled       string      `json:"is_enabled"`
	ActivateFlag    string      `json:"activate_flag"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdateAt        time.Time   `json:"update_at"`
	DeviceNumber    string      `json:"device_number"`
	ProductID       interface{} `json:"product_id"`
	ParentID        string      `json:"parent_id"`
	Protocol        interface{} `json:"protocol"`
	Label           string      `json:"label"`
	Location        interface{} `json:"location"`
	SubDeviceAddr   string      `json:"sub_device_addr"`
	CurrentVersion  interface{} `json:"current_version"`
	AdditionalInfo  string      `json:"additional_info"`
	ProtocolConfig  string      `json:"protocol_config"`
	Remark1         interface{} `json:"remark1"`
	Remark2         interface{} `json:"remark2"`
	Remark3         interface{} `json:"remark3"`
	DeviceConfigID  string      `json:"device_config_id"`
	BatchNumber     interface{} `json:"batch_number"`
	ActivateAt      interface{} `json:"activate_at"`
	IsOnline        int         `json:"is_online"`
	AccessWay       string      `json:"access_way"`
	Description     interface{} `json:"description"`
	ServiceAccessID string      `json:"service_access_id"`
}

type DeviceConfig struct {
	ID                     string `json:"id"`
	Voucher                string `json:"voucher"`
	DeviceType             string `json:"device_type"`
	ProtocolType           string `json:"protocol_type"`
	Config                 any    `json:"config"`
	ProtocolConfigTemplate any    `json:"protocol_config_template"`
	SubDevices             any    `json:"sub_devices"`
}
