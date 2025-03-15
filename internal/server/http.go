package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/docs"
	"github.com/quydmfl/niveau-test/internal/handler"
	"github.com/quydmfl/niveau-test/internal/middleware"
	"github.com/quydmfl/niveau-test/pkg/jwt"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/quydmfl/niveau-test/pkg/server/http"
	"github.com/quydmfl/niveau-test/pkg/validators"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	categoryHandler *handler.CategoryHandler,
	supplierHandler *handler.SupplierHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)

	// register custom validators
	validators.RegisterCustomValidators()

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			"message": "welcome to api backend!",
		})
	})

	s.GET("/health", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			"status": "healthy",
		})
	})

	api := s.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// users //
			auth := v1.Group("/")
			{
				auth.POST("/auth/register", userHandler.Register)
				auth.POST("/auth/login", userHandler.Login)
			}

			// products //
			products := v1.Group("/products").Use(middleware.StrictAuth(jwt, logger))
			{
				products.GET("/export/:format/:id", productHandler.ExportProducts)
				products.GET("/distance/ip/:city", productHandler.GetDistanceBetweenIPAndCity)
				products.GET("/", productHandler.GetProducts)
				products.GET("/:id", productHandler.GetProductDetail)
				products.POST("/", productHandler.CreateProduct)
				products.PUT("/:id", productHandler.UpdateProduct)
				products.DELETE("/:id", productHandler.DeleteProduct)
			}

			// categories //
			categories := v1.Group("/categories").Use(middleware.StrictAuth(jwt, logger))
			{
				categories.GET("/", categoryHandler.GetCategories)
				categories.GET("/:id", categoryHandler.GetCategoryDetail)
				categories.POST("/", categoryHandler.CreateCategory)
			}

			// supplier //
			supplier := v1.Group("/supplier").Use(middleware.StrictAuth(jwt, logger))
			{
				supplier.GET("/", supplierHandler.GetSuppliers)
				supplier.GET("/:id", supplierHandler.GetSupplierDetail)
				supplier.POST("/", supplierHandler.CreateSupplier)
			}

			// statistics //
			statistics := v1.Group("/statistics").Use(middleware.StrictAuth(jwt, logger))
			{
				statistics.GET("/products-per-category", productHandler.GetProductsPerCategory)
				statistics.GET("/products-per-supplier", productHandler.GetProductsPerSupplier)
			}
		}
	}

	return s
}
