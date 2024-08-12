package tpapi

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

// GetServiceList 获取服务接入点列表和设备列表
func GetServiceList(base string, identifier string) ([]Service, error) {
	resp, err := newRequest(base).
		SetBody(fmt.Sprintf(`{"service_identifier":"%s"}`, identifier)).
		SetResult(&ServiceListResult{}).
		Post("api/v1/plugin/service/access/list")

	if err != nil {
		return nil, err
	}

	data := resp.Result().(*ServiceListResult)
	if data.Code != 200 {
		return nil, errors.New(fmt.Sprintf("code=%d,msg=%s", data.Code, data.Message))
	}

	return data.Data, nil
}

// GetServiceById 获取服务详情包括设备列表
func GetServiceById(base string, id string) (*Service, error) {
	resp, err := newRequest(base).
		SetBody(fmt.Sprintf(`{"service_access_id":"%s"}`, id)).
		SetResult(&ServiceDetailResult{}).
		Post("api/v1/plugin/service/access")

	if err != nil {
		return nil, err
	}

	data := resp.Result().(*ServiceDetailResult)
	if data.Code != 200 {
		return nil, errors.New(data.Message)
	}

	return &data.Data, nil
}

// GetDeviceById 通过设备ID获取配置
func GetDeviceById(base string, id string) (*DeviceConfig, error) {
	return getDeviceConfig(base, "device_id", id)
}

// GetDeviceByNumber 通过设备编号获取配置
func GetDeviceByNumber(base string, number string) (*DeviceConfig, error) {
	return getDeviceConfig(base, "device_number", number)
}

// GetDeviceByVoucher 通过设备凭证获取配置
func GetDeviceByVoucher(base string, voucher string) (*DeviceConfig, error) {
	return getDeviceConfig(base, "voucher", voucher)
}

// SendHeartbeat 发送心跳
func SendHeartbeat(base string, identifier string) error {
	resp, err := newRequest(base).
		SetBody(fmt.Sprintf(`{"service_identifier":"%s"}`, identifier)).
		SetResult(&ApiResult{}).
		Post("api/v1/plugin/heartbeat")
	if err != nil {
		return err
	}

	data := resp.Result().(*ApiResult)
	if data.Code != 200 {
		return errors.New(data.Message)
	}
	return nil
}

// 获取设备配置信息
func getDeviceConfig(base string, auth string, value string) (*DeviceConfig, error) {
	request := newRequest(base)

	switch auth {
	case "device_id":
		request.SetBody(fmt.Sprintf(`{"device_id":"%s"}`, value))
	case "device_number":
		request.SetBody(fmt.Sprintf(`{"device_number":"%s"}`, value))
	case "voucher":
		request.SetBody(fmt.Sprintf(`{"voucher":"%s"}`, value))
	}

	resp, err := request.SetResult(&DeviceConfigResult{}).
		Post("api/v1/plugin/device/config")
	if err != nil {
		return nil, err
	}

	data := resp.Result().(*DeviceConfigResult)
	if data.Code != 200 {
		return nil, errors.New(data.Message)
	}

	return &data.Data, nil
}

func newRequest(base string) *resty.Request {
	client := resty.New()
	client.SetBaseURL(base)
	return client.R().SetHeader("Content-Type", "application/json")
}
