package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pplx2api/config"
	"pplx2api/logger"
	"strings"
	"sync"
	"time"
)

/**
 * 返回两个整数中的较小值
 */
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/**
 * 代理信息结构
 */
type ProxyInfo struct {
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ErrorCount int       `json:"error_count"`
}

/**
 * 代理池管理器
 * 负责从代理池API获取代理并管理代理的使用
 */
type ProxyPool struct {
	proxies []ProxyInfo
	mutex   sync.RWMutex
	index   int
}

/**
 * 代理池API响应结构
 */
type ProxyResponse struct {
	Proxy string `json:"proxy"`
	Error string `json:"error,omitempty"`
}

var (
	poolInstance *ProxyPool
	poolOnce     sync.Once
)

/**
 * 获取代理池单例实例
 */
func GetProxyPool() *ProxyPool {
	poolOnce.Do(func() {
		poolInstance = &ProxyPool{
			proxies: make([]ProxyInfo, 0),
			mutex:   sync.RWMutex{},
			index:   0,
		}
		// 如果启用了代理池，则初始化代理
		if config.ConfigInstance.EnableProxyPool {
			poolInstance.initializeProxies()
		}
	})
	return poolInstance
}

/**
 * 初始化代理池
 * 从代理池API获取指定数量的代理
 */
func (p *ProxyPool) initializeProxies() {
	if config.ConfigInstance.ProxyPoolAPI == "" {
		logger.Error("Proxy pool API not configured")
		return
	}

	logger.Info(fmt.Sprintf("Initializing proxy pool with %d proxies", config.ConfigInstance.ProxyPoolSize))

	// 并发获取代理以提高效率
	var wg sync.WaitGroup
	proxyChan := make(chan string, config.ConfigInstance.ProxyPoolSize)

	for i := 0; i < config.ConfigInstance.ProxyPoolSize; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			proxy, err := p.fetchProxyFromAPI()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to fetch proxy %d: %v", index, err))
				return
			}
			proxyChan <- proxy
		}(i)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(proxyChan)
	}()

	// 收集获取到的代理
	for proxy := range proxyChan {
		p.mutex.Lock()
		proxyInfo := ProxyInfo{
			URL:        proxy,
			CreatedAt:  time.Now(),
			ErrorCount: 0,
		}
		p.proxies = append(p.proxies, proxyInfo)
		p.mutex.Unlock()
	}

	logger.Info(fmt.Sprintf("Successfully initialized proxy pool with %d proxies", len(p.proxies)))
}

/**
 * 从代理池API获取单个代理
 */
func (p *ProxyPool) fetchProxyFromAPI() (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(config.ConfigInstance.ProxyPoolAPI)
	if err != nil {
		return "", fmt.Errorf("failed to request proxy API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("proxy API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应是否为空
	if len(body) == 0 {
		return "", fmt.Errorf("empty response from proxy API")
	}

	bodyStr := strings.TrimSpace(string(body))

	// 检查响应是否是纯文本代理URL格式
	if strings.HasPrefix(bodyStr, "http://") || strings.HasPrefix(bodyStr, "https://") {
		// 直接返回代理URL
		return bodyStr, nil
	}

	// 尝试解析JSON格式
	var proxyResp ProxyResponse
	if err := json.Unmarshal(body, &proxyResp); err != nil {
		return "", fmt.Errorf("failed to parse proxy response as JSON (body: %s): %w", bodyStr, err)
	}

	if proxyResp.Error != "" {
		return "", fmt.Errorf("proxy API error: %s", proxyResp.Error)
	}

	if proxyResp.Proxy == "" {
		return "", fmt.Errorf("empty proxy returned from API")
	}

	logger.Info(fmt.Sprintf("Fetched proxy: %s", proxyResp.Proxy))
	return proxyResp.Proxy, nil
}

/**
 * 获取下一个可用的代理
 * 使用轮询方式分配代理
 */
func (p *ProxyPool) GetNextProxy() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.proxies) == 0 {
		logger.Error("No proxies available in pool")
		return ""
	}

	proxy := p.proxies[p.index]
	p.index = (p.index + 1) % len(p.proxies)

	return proxy.URL
}

/**
 * 随机获取一个代理
 */
func (p *ProxyPool) GetRandomProxy() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if len(p.proxies) == 0 {
		logger.Error("No proxies available in pool")
		return ""
	}

	index := rand.Intn(len(p.proxies))
	return p.proxies[index].URL
}

/**
 * 获取代理池大小
 */
func (p *ProxyPool) Size() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return len(p.proxies)
}

/**
 * 刷新代理池
 * 重新从API获取所有代理
 */
func (p *ProxyPool) RefreshPool() {
	logger.Info("Refreshing proxy pool")

	p.mutex.Lock()
	p.proxies = make([]ProxyInfo, 0)
	p.index = 0
	p.mutex.Unlock()

	p.initializeProxies()
}

/**
 * 移除失效的代理
 */
func (p *ProxyPool) RemoveProxy(proxy string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, existingProxy := range p.proxies {
		if existingProxy.URL == proxy {
			p.proxies = append(p.proxies[:i], p.proxies[i+1:]...)
			logger.Info(fmt.Sprintf("Removed proxy: %s", proxy))
			break
		}
	}

	// 调整索引
	if p.index >= len(p.proxies) && len(p.proxies) > 0 {
		p.index = 0
	}
}

/**
 * 添加新代理到池中
 */
func (p *ProxyPool) AddProxy(proxy string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 检查代理是否已存在
	for _, existingProxy := range p.proxies {
		if existingProxy.URL == proxy {
			return
		}
	}

	proxyInfo := ProxyInfo{
		URL:        proxy,
		CreatedAt:  time.Now(),
		ErrorCount: 0,
	}
	p.proxies = append(p.proxies, proxyInfo)
	logger.Info(fmt.Sprintf("Added proxy to pool: %s", proxy))
}

/**
 * 处理代理错误，如果是407错误则删除代理
 */
func (p *ProxyPool) HandleProxyError(proxyURL string, statusCode int) {
	if statusCode == 407 {
		logger.Error(fmt.Sprintf("Proxy authentication failed (407), removing proxy: %s", proxyURL))
		p.RemoveProxy(proxyURL)
	} else {
		// 增加错误计数
		p.mutex.Lock()
		defer p.mutex.Unlock()

		for i, proxyInfo := range p.proxies {
			if proxyInfo.URL == proxyURL {
				p.proxies[i].ErrorCount++
				logger.Info(fmt.Sprintf("Proxy error count increased for %s: %d", proxyURL, p.proxies[i].ErrorCount))

				// 如果错误次数过多，也删除代理
				if p.proxies[i].ErrorCount >= 5 {
					logger.Error(fmt.Sprintf("Proxy error count exceeded limit, removing proxy: %s", proxyURL))
					p.proxies = append(p.proxies[:i], p.proxies[i+1:]...)
					// 调整索引
					if p.index >= len(p.proxies) && len(p.proxies) > 0 {
						p.index = 0
					}
				}
				break
			}
		}
	}
}

/**
 * 获取代理详细信息
 */
func (p *ProxyPool) GetProxyInfo() []ProxyInfo {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	proxies := make([]ProxyInfo, len(p.proxies))
	copy(proxies, p.proxies)
	return proxies
}

/**
 * 获取所有代理列表
 */
func (p *ProxyPool) GetAllProxies() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	proxies := make([]string, len(p.proxies))
	for i, proxyInfo := range p.proxies {
		proxies[i] = proxyInfo.URL
	}
	return proxies
}

/**
 * 检查是否需要轮换代理池
 * 如果代理池中的代理创建时间超过6小时，则需要轮换
 */
func (p *ProxyPool) ShouldRotate() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if len(p.proxies) == 0 {
		return true
	}

	// 检查最老的代理是否超过6小时
	oldestTime := time.Now()
	for _, proxyInfo := range p.proxies {
		if proxyInfo.CreatedAt.Before(oldestTime) {
			oldestTime = proxyInfo.CreatedAt
		}
	}

	return time.Since(oldestTime) > 6*time.Hour
}
