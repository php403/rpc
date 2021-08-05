package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/php403/im/internal/logic/conf"
	"github.com/php403/im/internal/logic/service"
	"github.com/php403/im/pkg/log"
	"net/http"
)


func NewHTTPServer(conf *conf.Server, logger log.Logger,service *service.UserService) *http.Server {
	fmt.Println(conf.Http)
	s := &http.Server{
		Addr : conf.Http.Addr,
		Handler: NewRouter(service),
		ReadTimeout: conf.Http.ReadTimeout,
		WriteTimeout: conf.Http.WriteTimeout,
		MaxHeaderBytes: conf.Http.MaxHeaderBytes,
	}
	return s
}

func NewRouter(service *service.UserService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("api")
	{
		api.POST("/auth",service.Auth)
	}
	return r
}




