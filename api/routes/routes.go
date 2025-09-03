package routes

import (
	"URLShortener/api/controllers"
	"URLShortener/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, uc *controllers.UserController, lc *controllers.LinkController, ac *controllers.AuthController) {
	router.Use(gin.Recovery())           // 500 error handler
	router.Use(gin.Logger())             // Logger
	router.Use(middleware.RateLimiter()) // Rate limiter

	apiv1 := router.Group("/api/v1")

	// public endpoints for auth
	apiv1.POST("/login", ac.Login)
	apiv1.POST("/register", ac.Register)

	// add JWT authentication for login-dependent endpoints
	protectedApiv1 := apiv1.Use(middleware.AuthMiddleware())

	// /logout endpoint
	protectedApiv1.POST("/logout", ac.Logout) // logout requires the user to be logged in

	// /user endpoint
	protectedApiv1.GET("/user", uc.GetAllUsers)
	protectedApiv1.GET("/user/:id", uc.GetUser)
	protectedApiv1.PUT("/user/:id", uc.UpdateUser)
	protectedApiv1.POST("/user/:id", uc.ReactivateUser)
	protectedApiv1.DELETE("/user/:id", uc.DeactivateUser)
	protectedApiv1.DELETE("/user/:id/force", uc.DeleteUser)

	// /link endpoint
	protectedApiv1.GET("/short_links", lc.GetAllLinksByUser)
	protectedApiv1.GET("/short_links/:id", lc.GetLink)
	protectedApiv1.GET("/short_link/:code", lc.GetFullLink)
	protectedApiv1.POST("/short_links", lc.CreateLink)
	protectedApiv1.PUT("/short_links", lc.UpdateLink)
	protectedApiv1.DELETE("/short_links/:id", lc.DeactivateLink)
	protectedApiv1.DELETE("/short_links/:id/force", lc.DeleteLink)
}
