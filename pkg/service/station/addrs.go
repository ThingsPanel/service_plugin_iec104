package station

// IEC104规约的地址范围主要涉及到遥信、‌遥测、‌遥控等信息的基地址，‌这些地址范围在1997年和2002年的版本中有所不同。‌在2002版的IEC104规约中，
// ‌遥信信息的基地址范围从1H到4000H，
// ‌遥测信息的基地址范围从4001H到5000H，‌
// 遥控信息的基地址范围从6001H到6100H
// ‌设点信息的基地址范围从6201H到6400H，‌
// 电度信息的基地址范围从6401H到6600H

type DeviceType string

const (
	YaoXin DeviceType = "yx" // 遥信信息
	YaoCe  DeviceType = "yc" // 遥测信息
	DianDu DeviceType = "dd" // 电度信息
)

func GetDeviceType(addr uint16) DeviceType {
	if addr <= 0x400 {
		return YaoXin
	} else if addr >= 0x401 && addr <= 0x500 {
		return YaoCe
	} else if addr >= 0x6401 && addr <= 0x6600 {
		return DianDu
	}
	return "未知设备类型"
}
