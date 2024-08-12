package plugin

import (
	"context"
	"iec104-slave/pkg/plugin/voucher"
)

// Provider 插件外部服务提供
type Provider interface {
	// OnRequestDevices 获取设备列表
	OnRequestDevices(ctx context.Context, voucher *voucher.Service, page int, size int) ([]Device, int, error)

	// OnEventNotify 事件通知
	OnEventNotify(ctx context.Context)

	// OnDisconnectDevice 断开设备链接
	OnDisconnectDevice(ctx context.Context, deviceId string) error
}
