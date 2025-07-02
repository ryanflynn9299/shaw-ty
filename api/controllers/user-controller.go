package controllers

import (
	"URLShortener/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	userService services.UserService
}

// TODO: Finish implementing this file

// NewUserController initializes a new UserController
func NewUserController(userService *services.UserService) UserController {
	return UserController{
		userService: *userService,
	}
}

// GetUser gets a user with the given ID
func (uctlr *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := uctlr.userService.GetUserById(id)
	if err != nil {
		log.Println("Error getting user:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers retrieves all users
func (uctlr *UserController) GetAllUsers(c *gin.Context) {
	users, err := uctlr.userService.GetAllUsers(c)
	if err != nil {
		log.Println("Error getting all users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser updates any provided fields. is idempotent
func (uctlr *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var userData struct {
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = uctlr.userService.UpdateUserById(id, userData.Password, userData.FirstName, userData.LastName, userData.Email)
	if err != nil {
		log.Println("Error updating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeactivateUser soft-deletes the user by deactivating them, preferred over delete
func (uctlr *UserController) DeactivateUser(c *gin.Context) {
	// Note: This would require adding a deactivation method to the UserService interface
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// ReactivateUser enables a previously soft-deleted user (admin-only workflow)
func (uctlr *UserController) ReactivateUser(c *gin.Context) {
	// Note: This would require adding a reactivation method to the UserService interface
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// DeleteUser hard-deletes the user from the DB and erases all related data. Used for privacy compliance,
// but not a common workflow.
func (uctlr *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = uctlr.userService.DeleteUserById(id)
	if err != nil {
		log.Println("Error deleting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
