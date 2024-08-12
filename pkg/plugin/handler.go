package plugin

import (
	"context"
	"iec104-slave/pkg/plugin/server"
	"iec104-slave/pkg/plugin/voucher"
)

type handler struct {
	provider Provider
}

func newHandler(provider Provider) *handler {
	return &handler{
		provider: provider,
	}
}

func (s *handler) OnDeviceListRequest(ctx context.Context, req *server.DevicesRequest) ([]server.Device, int, error) {
	voc, err := voucher.NewService().FromJson(req.Voucher)
	if err != nil {
		return nil, 0, err
	}

	list, total, err := s.provider.OnRequestDevices(ctx, voc, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]server.Device, 0, len(list))
	for _, item := range list {
		result = append(result, server.Device{
			DeviceNumber: item.DeviceNumber(),
			DeviceName:   item.DeviceName(),
			Description:  item.Description(),
		})
	}

	return result, total, nil
}

func (s *handler) OnNotifyEvent(ctx context.Context) {
	s.provider.OnEventNotify(ctx)
}

func (s *handler) OnDisconnectDevice(ctx context.Context, deviceId string) error {
	return s.provider.OnDisconnectDevice(ctx, deviceId)
}
