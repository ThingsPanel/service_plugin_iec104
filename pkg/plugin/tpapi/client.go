package tpapi

type Client struct {
	BaseUrl    string
	Identifier string
}

func NewClient(baseUrl string, identifier string) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		Identifier: identifier,
	}
}

// GetServiceList 获取服务接入点列表和设备列表
func (c *Client) GetServiceList() ([]Service, error) {
	return GetServiceList(c.BaseUrl, c.Identifier)
}

// GetServiceById 获取服务详情包括设备列表
func (c *Client) GetServiceById(id string) (*Service, error) {
	return GetServiceById(c.BaseUrl, id)
}

// GetDeviceById 通过设备ID获取配置
func (c *Client) GetDeviceById(id string) (*DeviceConfig, error) {
	return GetDeviceById(c.BaseUrl, id)
}

// GetDeviceByNumber 通过设备编号获取配置
func (c *Client) GetDeviceByNumber(number string) (*DeviceConfig, error) {
	return GetDeviceByNumber(c.BaseUrl, number)
}

// GetDeviceByVoucher 通过设备凭证获取配置
func (c *Client) GetDeviceByVoucher(voucher string) (*DeviceConfig, error) {
	return GetDeviceByVoucher(c.BaseUrl, voucher)
}

// SendHeartbeat 发送心跳
func (c *Client) SendHeartbeat() error {
	return SendHeartbeat(c.BaseUrl, c.Identifier)
}
