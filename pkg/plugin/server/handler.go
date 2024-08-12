package server

import "context"

type Handler interface {
	// OnDeviceListRequest 设备列表请求
	OnDeviceListRequest(ctx context.Context, req *DevicesRequest) ([]Device, int, error)

	// OnNotifyEvent 服务配置变化通知
	OnNotifyEvent(ctx context.Context)

	// OnDisconnectDevice 设备断开刷新通知
	OnDisconnectDevice(ctx context.Context, deviceId string) error
}
