package routes

import (
	"URLShortener/api/controllers"
	"URLShortener/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, uc *controllers.UserController, lc *controllers.LinkController, ac *controllers.AuthController) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

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
	protectedApiv1.GET("/short_link/:id", lc.GetLink)
	protectedApiv1.GET("/short_link", lc.GetFullLink)
	protectedApiv1.POST("/short_link", lc.CreateLink)
	protectedApiv1.PUT("/short_link", lc.UpdateLink)
	protectedApiv1.DELETE("/short_link/:id", lc.DeactivateLink)
	protectedApiv1.DELETE("/short_link/:id/force", lc.DeleteLink)

}
