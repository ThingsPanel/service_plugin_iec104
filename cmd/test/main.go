package main

import (
	"fmt"
	"time"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type myClient struct{}

func main() {
	var err error

	option := cs104.NewOption()
	option.SetConfig(cs104.Config{
		ConnectTimeout0:   30 * time.Second,
		SendUnAckLimitK:   12,
		SendUnAckTimeout1: 15 * time.Second,
		RecvUnAckLimitW:   2,
		RecvUnAckTimeout2: 10 * time.Second,
		IdleTimeout3:      20 * time.Second,
	})

	if err = option.AddRemoteServer("127.0.0.1:2404"); err != nil {
		panic(err)
	}

	mycli := &myClient{}

	client := cs104.NewClient(mycli, option)

	client.LogMode(true)

	client.SetOnConnectHandler(func(c *cs104.Client) {
		c.SendStartDt() // 发送startDt激活指令
	})

	client.SetActiveHandler(func(c *cs104.Client) {
		err := c.InterrogationCmd(
			asdu.CauseOfTransmission{
				Cause: asdu.Activation,
			},
			asdu.CommonAddr(1),
			asdu.QOIStation,
		)
		if err != nil {
			panic(err)
		}
	})

	err = client.Start()
	if err != nil {
		panic(fmt.Errorf("Failed to connect. error:%v\n", err))
	}

	for {
		time.Sleep(time.Second * 100)
	}

}

func (myClient) InterrogationHandler(c asdu.Connect, a *asdu.ASDU) error {
	fmt.Println("InterrogationHandler", a.CommonAddr)
	for _, item := range a.GetSinglePoint() {
		fmt.Println(item.Ioa)
	}
	return nil
}

func (myClient) CounterInterrogationHandler(c asdu.Connect, a *asdu.ASDU) error {
	fmt.Println("CounterInterrogationHandler", a)
	return nil
}

func (myClient) ReadHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) TestCommandHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) ClockSyncHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) ResetProcessHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (myClient) ASDUHandler(c asdu.Connect, a *asdu.ASDU) error {
	fmt.Println("ASDUHandler", a.Type, a.CommonAddr, a.Coa)
	return nil
}
