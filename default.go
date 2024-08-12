package iec104_master

import (
	"github.com/gin-gonic/gin"
	"iec104-slave/pkg/config"
	"iec104-slave/pkg/logger"
)

var Default *Server

func Init(cfg ...string) error {
	if err := config.Load(cfg...); err != nil {
		return err
	}
	logger.Init(config.Default.Logs)

	if !config.Default.Logs.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	Default = NewServer(config.Default)

	return nil
}

func Run() error {
	return Default.Run()
}

func Close() error {
	return nil
}
