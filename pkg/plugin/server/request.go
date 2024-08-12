package server

// FormRequest 服务配置表单请求
type FormRequest struct {
	ProtocolType string `form:"protocol_type" binding:"required"`
	DeviceType   string `form:"device_type"`
	FormType     string `form:"form_type" binding:"required"`
}

// DevicesRequest 服务设备列表请求
type DevicesRequest struct {
	ServiceIdentifier string `form:"service_identifier" json:"service_identifier"`
	Voucher           string `form:"voucher" json:"voucher"`
	Page              int    `form:"page" json:"page"`
	PageSize          int    `form:"page_size" json:"page_size"`
}
