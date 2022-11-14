package http

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/a-line-cli/pkg/service"
)

type HttpServer struct {
	jobService service.IJobService
}

func NewHttpService(jobService service.IJobService) *HttpServer {
	return &HttpServer{
		jobService: jobService,
	}
}

func (h *HttpServer) StartHttpServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		//输出json结果给调用方
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080") // listen and serve on
}
