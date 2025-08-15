package router

import (
	"pplx2api/middleware"
	"pplx2api/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Apply middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	// Health check endpoint
	r.GET("/health", service.HealthCheckHandler)

	// Chat completions endpoint (OpenAI-compatible)
	r.POST("/v1/chat/completions", service.ChatCompletionsHandler)
	r.GET("/v1/models", service.ModelsHandler)

	// Proxy pool management endpoints
	r.GET("/proxy/status", service.ProxyPoolStatusHandler)
	r.POST("/proxy/refresh", service.RefreshProxyPoolHandler)
	// HuggingFace compatible routes
	hfRouter := r.Group("/hf")
	{
		v1Router := hfRouter.Group("/v1")
		{
			v1Router.POST("/chat/completions", service.ChatCompletionsHandler)
			v1Router.GET("/models", service.ModelsHandler)
		}
	}
}
