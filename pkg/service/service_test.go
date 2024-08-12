package service

import "testing"

func TestNewService(t *testing.T) {
	svc := NewService(Config{
		Remote: "127.0.0.1:2404",
	})

	svc.Run()
}
