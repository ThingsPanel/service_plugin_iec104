package tpapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func Test_GetServiceList(t *testing.T) {
	res, err := GetServiceList("http://water.thingspanel.cn", "IEC104")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func Test_GetServiceDetail(t *testing.T) {
	res, err := GetServiceById("http://c.thingspanel.cn", "f8df6ae6-1a59-d623-9b2e-70323f822caa")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func Test_GetDeviceConfig(t *testing.T) {
	client := resty.New()
	client.SetBaseURL("http://c.thingspanel.cn")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"device_number":"10000"}`).
		Post("api/v1/plugin/device/config")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.String())
}

func Test_GetDeviceByNumber(t *testing.T) {
	resp, err := GetDeviceByNumber("http://c.thingspanel.cn", "10000")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func Test_SendHeartbeat(t *testing.T) {
	err := SendHeartbeat("http://c.thingspanel.cn", "test")
	if err != nil {
		t.Fatal(err)
	}
}
