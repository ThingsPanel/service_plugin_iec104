package service

type IpAddr struct {
	addr string
}

func (i *IpAddr) Addr() string {
	return i.addr
}
