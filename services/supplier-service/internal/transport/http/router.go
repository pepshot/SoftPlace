package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/pepshot/SoftPlace/services/supplier-service/internal/transport/http/middleware"
	"github.com/pepshot/SoftPlace/shared/auth"
)

func NewRouter(handler *Handler, jwtManager *auth.JWTManager) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", handler.Register)
		v1.POST("/auth/login", handler.Login)
	}

	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		protected.GET("/profile", handler.Profile)

		protected.GET("/modules", handler.ListModules)
		protected.POST("/modules", handler.CreateModule)
		protected.GET("/modules/:id", handler.GetModule)
		protected.PUT("/modules/:id", handler.UpdateModule)
		protected.DELETE("/modules/:id", handler.DeleteModule)

		protected.GET("/materials", handler.ListMaterials)
		protected.POST("/materials", handler.CreateMaterial)
		protected.GET("/materials/:id", handler.GetMaterial)
		protected.PUT("/materials/:id", handler.UpdateMaterial)
		protected.DELETE("/materials/:id", handler.DeleteMaterial)

		protected.GET("/supplies", handler.ListSupplies)
		protected.POST("/supplies", handler.CreateSupply)
		protected.GET("/supplies/:id", handler.GetSupply)
		protected.PUT("/supplies/:id", handler.UpdateSupply)
		protected.DELETE("/supplies/:id", handler.DeleteSupply)
	}

	return router
}
