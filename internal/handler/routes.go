package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

// SetupRoutes registra todas as rotas da sua API
func SetupRoutes(userHandler *UserHandler, characterHandler *CharacterHandler) *gin.Engine {
	router := gin.Default()

	// Rota do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas gerais
	router.GET("/health", HealthCheckHandler)

	// Rotas de Usuários
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUserById)
	router.DELETE("/users/:id", userHandler.DeleteUserById)
	router.PATCH("/users/:id", userHandler.UpdateUser)

	// Rotas de Personagens
	router.POST("/characters", characterHandler.CreateCharacter)
	router.GET("/characters/:id", characterHandler.GetCharacterById)
	router.DELETE("/characters/:id", characterHandler.DeleteCharacterById)
	router.PATCH("/characters/:id", characterHandler.UpdateCharacter)

	return router
}

func Setup() *gin.Engine {
	// Inicialização de serviços e handlers
	// User
	userService := services.NewUserService(database.DB)
	userHandler := NewUserHandler(userService)

	// Character
	characterService := services.NewCharacterService(database.DB)
	characterHandler := NewCharacterHandler(characterService)

	return SetupRoutes(userHandler, characterHandler)
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
