package tpapi

import (
	"iec104-slave/pkg/plugin/voucher"
	"time"
)

type Service struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Voucher         string    `json:"voucher"`
	Devices         []Device  `json:"devices"`
	ServicePluginId string    `json:"service_plugin_id"`
	TenantId        string    `json:"tenant_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdateAt        time.Time `json:"update_at"`
}

func (s Service) Vouchers() (*voucher.Service, error) {
	return voucher.NewService().FromJson(s.Voucher)
}
