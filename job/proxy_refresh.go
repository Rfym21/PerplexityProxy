package job

import (
	"fmt"
	"pplx2api/config"
	"pplx2api/logger"
	"pplx2api/proxy"
	"sync"
	"time"
)

/**
 * 代理池刷新器
 * 定期刷新代理池中的代理
 */
type ProxyRefresher struct {
	interval  time.Duration
	stopChan  chan struct{}
	isRunning bool
	mutex     sync.Mutex
}

var (
	proxyRefresherInstance *ProxyRefresher
	proxyRefresherOnce     sync.Once
)

/**
 * 获取代理池刷新器单例实例
 */
func GetProxyRefresher(interval time.Duration) *ProxyRefresher {
	proxyRefresherOnce.Do(func() {
		proxyRefresherInstance = &ProxyRefresher{
			interval:  interval,
			stopChan:  make(chan struct{}),
			isRunning: false,
		}
	})
	return proxyRefresherInstance
}

/**
 * 启动代理池刷新器
 */
func (pr *ProxyRefresher) Start() {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if pr.isRunning {
		logger.Info("Proxy refresher is already running")
		return
	}

	if !config.ConfigInstance.EnableProxyPool {
		logger.Info("Proxy pool is disabled, proxy refresher will not start")
		return
	}

	pr.isRunning = true
	logger.Info("Starting proxy refresher")

	go pr.run()
}

/**
 * 停止代理池刷新器
 */
func (pr *ProxyRefresher) Stop() {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if !pr.isRunning {
		return
	}

	logger.Info("Stopping proxy refresher")
	close(pr.stopChan)
	pr.isRunning = false
}

/**
 * 运行代理池刷新器
 */
func (pr *ProxyRefresher) run() {
	ticker := time.NewTicker(pr.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pr.refreshProxies()
		case <-pr.stopChan:
			logger.Info("Proxy refresher stopped")
			return
		}
	}
}

/**
 * 刷新代理池
 * 检查是否需要轮换代理，如果需要则轮换一批新的代理
 */
func (pr *ProxyRefresher) refreshProxies() {
	logger.Info("Checking proxy pool rotation")

	proxyPool := proxy.GetProxyPool()

	// 检查是否需要轮换
	if proxyPool.ShouldRotate() {
		logger.Info("Starting proxy pool rotation")
		oldSize := proxyPool.Size()

		// 轮换代理池
		proxyPool.RefreshPool()

		newSize := proxyPool.Size()
		logger.Info("Proxy pool rotation completed")
		logger.Info(fmt.Sprintf("Old proxy count: %d, New proxy count: %d", oldSize, newSize))
	} else {
		logger.Info("Proxy pool rotation not needed, proxies are still fresh")
	}
}

/**
 * 检查代理池状态
 */
func (pr *ProxyRefresher) IsRunning() bool {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	return pr.isRunning
}
