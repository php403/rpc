// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/php403/im/internal/logic/biz"
	"github.com/php403/im/internal/logic/conf"
	"github.com/php403/im/internal/logic/data"
	"github.com/php403/im/internal/logic/server"
	"github.com/php403/im/internal/logic/service"
	"github.com/php403/im/pkg/log"
)

func initApp(*conf.Server, *conf.Data, log.Logger) (*App, func(), error) {
		panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
