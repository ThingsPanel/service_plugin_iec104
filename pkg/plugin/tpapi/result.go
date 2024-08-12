package tpapi

type ApiResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ServiceListResult struct {
	ApiResult
	Data []Service `json:"data"`
}

type ServiceDetailResult struct {
	ApiResult
	Data Service `json:"data"`
}

type DeviceConfigResult struct {
	ApiResult
	Data DeviceConfig `json:"data"`
}
