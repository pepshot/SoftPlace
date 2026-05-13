package httpserver

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/pepshot/SoftPlace/services/customer-service/internal/transport/http/middleware"
	"github.com/pepshot/SoftPlace/shared/auth"
)

func NewRouter(handler *Handler, jwtManager *auth.JWTManager) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", handler.Register)
		v1.POST("/auth/login", handler.Login)
	}

	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		protected.GET("/profile", handler.Profile)

		protected.GET("/furniture", handler.ListFurniture)
		protected.POST("/furniture", handler.CreateFurniture)
		protected.GET("/furniture/:id", handler.GetFurniture)
		protected.PUT("/furniture/:id", handler.UpdateFurniture)
		protected.DELETE("/furniture/:id", handler.DeleteFurniture)

		protected.GET("/garnitures", handler.ListGarnitures)
		protected.POST("/garnitures", handler.CreateGarniture)
		protected.GET("/garnitures/:id", handler.GetGarniture)
		protected.PUT("/garnitures/:id", handler.UpdateGarniture)
		protected.DELETE("/garnitures/:id", handler.DeleteGarniture)

		protected.GET("/shipments", handler.ListShipments)
		protected.POST("/shipments", handler.CreateShipment)
		protected.GET("/shipments/:id", handler.GetShipment)
		protected.PUT("/shipments/:id", handler.UpdateShipment)
		protected.DELETE("/shipments/:id", handler.DeleteShipment)
	}

	return router
}
