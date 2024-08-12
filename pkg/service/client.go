package service

import (
	"github.com/thinkgos/go-iecp5/asdu"
	"iec104-slave/pkg/service/station"
	"log/slog"
)

type iecClient struct {
	station   *station.Stations
	onReceive OnIecReceive
}

func newIecClient(station *station.Stations, onReceive OnIecReceive) *iecClient {
	return &iecClient{
		station:   station,
		onReceive: onReceive,
	}
}

func (cli *iecClient) InterrogationHandler(c asdu.Connect, a *asdu.ASDU) error {
	return cli.HandleData(c, a)
}

func (cli *iecClient) CounterInterrogationHandler(conn asdu.Connect, data *asdu.ASDU) error {
	return cli.HandleData(conn, data)
}

func (cli *iecClient) ASDUHandler(conn asdu.Connect, data *asdu.ASDU) error {
	return cli.HandleData(conn, data)
}

func (cli *iecClient) TestCommandHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cli *iecClient) ClockSyncHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cli *iecClient) ResetProcessHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cli *iecClient) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cli *iecClient) ReadHandler(conn asdu.Connect, data *asdu.ASDU) error {
	return nil
}

func (cli *iecClient) HandleData(conn asdu.Connect, data *asdu.ASDU) error {
	addrType := station.GetDeviceType(uint16(data.CommonAddr))
	commonAddr := uint16(data.CommonAddr)
	stationItem := cli.station.Load(commonAddr)

	slog.Debug("station", "addr", data.CommonAddr, "type", addrType)

	parseIoaAndValue(data, func(ioa uint, value any) {
		stationItem.Devices.Add(ioa)
		cli.onReceive(commonAddr, ioa, value)
		slog.Debug("data", "ioa", ioa, "value", value)
	})

	slog.Debug("data", "station", addrType, "device", stationItem.Devices.Count())

	return nil
}

// 解析对象地址和值
func parseIoaAndValue(data *asdu.ASDU, cb func(ioa uint, value any)) {
	switch data.Type {
	case asdu.M_SP_NA_1, asdu.M_SP_TA_1:
		for _, item := range data.GetSinglePoint() {
			cb(uint(item.Ioa), item.Value)
		}
	case asdu.M_DP_NA_1, asdu.M_DP_TA_1:
		for _, item := range data.GetDoublePoint() {
			cb(uint(item.Ioa), item.Value)
		}
	case asdu.M_ST_NA_1, asdu.M_ST_TA_1:
		for _, item := range data.GetStepPosition() {
			cb(uint(item.Ioa), item.Value.Val)
		}
	case asdu.M_BO_NA_1, asdu.M_BO_TA_1:
		for _, item := range data.GetBitString32() {
			cb(uint(item.Ioa), item.Value)
		}
	case asdu.M_ME_NA_1, asdu.M_ME_TA_1:
		for _, item := range data.GetMeasuredValueNormal() {
			cb(uint(item.Ioa), item.Value)
		}
	case asdu.M_ME_NB_1, asdu.M_ME_TB_1:
		for _, item := range data.GetMeasuredValueScaled() {
			cb(uint(item.Ioa), item.Value)
		}
	}
}
