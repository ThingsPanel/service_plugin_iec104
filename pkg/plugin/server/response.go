package server

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiDevicesResponse struct {
	ApiResponse
	Data struct {
		List  []Device `json:"list"`
		Total int      `json:"total"`
	} `json:"data"`
}

func NewApiDevicesResponse(list []Device, total int) *ApiDevicesResponse {
	res := ApiDevicesResponse{}
	res.Code = 200
	res.Data.List = list[:]
	res.Data.Total = total
	return &res
}
