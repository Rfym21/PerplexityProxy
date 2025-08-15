package main

import (
	"pplx2api/config"
	"pplx2api/job"
	"pplx2api/proxy"
	"pplx2api/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Load configuration

	// 初始化代理池
	if config.ConfigInstance.EnableProxyPool {
		_ = proxy.GetProxyPool()
	}

	// Setup all routes
	router.SetupRoutes(r)
	// 创建会话更新器，设置更新间隔为24小时
	sessionUpdater := job.GetSessionUpdater(24 * time.Hour)

	// 启动会话更新器
	sessionUpdater.Start()
	defer sessionUpdater.Stop()

	// 创建代理池刷新器，设置更新间隔为6小时
	if config.ConfigInstance.EnableProxyPool {
		proxyRefresher := job.GetProxyRefresher(6 * time.Hour)
		proxyRefresher.Start()
		defer proxyRefresher.Stop()
	}

	// Run the server on 0.0.0.0:8080
	r.Run(config.ConfigInstance.Address)
}
