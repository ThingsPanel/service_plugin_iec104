package service

import "context"

type Config struct {
	ServiceId string
	Remote    string
	Handler   OnDataReceive
	Context   context.Context
}
