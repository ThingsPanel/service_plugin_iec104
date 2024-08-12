package server

import (
	"github.com/gin-gonic/gin"
	form2 "iec104-slave/pkg/plugin/form"
	"log/slog"
	"net/http"
)

type Server struct {
	http    *http.Server
	router  *gin.Engine
	handler Handler
}

func NewServer(handler Handler) *Server {
	router := gin.Default()
	server := &Server{
		router:  router,
		handler: handler,
		http:    &http.Server{},
	}
	server.init()
	return server
}

func (s *Server) init() {
	router := s.router.Group("/api/v1")
	// 获取表单配置
	router.GET("/form/config", s.handleFormConfig)
	// 获取设备列表
	router.GET("plugin/device/list", s.handleDeviceList)
	// 通知回调
	router.POST("notify/event", s.handleNotifyEvent)
	// 断开设备
	router.POST("device/disconnect", s.handleDeviceDisconnect)
}

func (s *Server) handleFormConfig(c *gin.Context) {
	var req FormRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switch req.FormType {
	case "CFG": // 配置表单
		c.Data(200, "application/json", []byte(form2.DataConfig))
	case "VCR": // 凭证表单
		c.Data(200, "application/json", []byte(form2.DeviceVoucher))
	case "SVCR": // 服务凭证表单
		c.Data(200, "application/json", []byte(form2.ServiceVoucher))
	default:
		c.Status(404)
	}
}

func (s *Server) handleDeviceList(c *gin.Context) {
	var req DevicesRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("device list request data", "err", err)
		c.JSON(http.StatusOK, &ApiResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	list, total, err := s.handler.OnDeviceListRequest(c, &req)
	if err != nil {
		slog.Warn("device list request result", "err", err)
		c.JSON(http.StatusOK, &ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, NewApiDevicesResponse(list, total))
}

func (s *Server) handleNotifyEvent(c *gin.Context) {
	// 如果message_type是1，message的内容为{"service_access_id":"xxxxx"}
	s.handler.OnNotifyEvent(c)
}

func (s *Server) handleDeviceDisconnect(c *gin.Context) {
	id, ok := c.GetPostForm("davice_id")
	if ok {
		s.handler.OnDisconnectDevice(c, id)
	}
}

func (s *Server) Run(addr string) error {
	s.http.Addr = addr
	s.http.Handler = s.router.Handler()
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	return s.http.Close()
}
