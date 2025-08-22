package controllers

import (
	"URLShortener/api/dto"
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
// TODO: migrate strings to yaml file
// TODO: add input sanitization and validation

// NewUserController initializes a new UserController
func NewUserController(userService *services.UserService) UserController {
	return UserController{
		userService: *userService,
	}
}

// GetUser gets a user with the given ID
func (uctlr *UserController) GetUser(c *gin.Context) {
	// TODO: prevent access to other users unless they are the owner, RBAP and OWASP compliance
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

	// TODO: return limited PII version: Admins don't need to see PII
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

	// TODO: defer dto conversion to service
	var userRequest dto.UpdateUserRequest
	support := c.Query("support")
	if support == "true" {
		// TODO: remove non-support updates
		err = c.ShouldBindBodyWithJSON(&userRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The request body provided is invalid, please check your request and try again."})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The request body provided is invalid, please check your request and try again."})
			return
		}
	}

	success, err := uctlr.userService.UpdateUserById(id, userRequest.Username, userRequest.Password, userRequest.FirstName, userRequest.LastName, userRequest.Email, userRequest.IsActive)
	if err != nil || success == 1 {
		log.Println("Error updating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user, please check your request and try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully."})
}

// DeactivateUser soft-deletes the user by deactivating them, preferred over delete
func (uctlr *UserController) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid user ID format. Please try again."})
		return
	}
	// TODO: add user exists check and guard

	success, err := uctlr.userService.DeactivateUserById(id)
	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Failed to deactivate user, please try again."})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deactivated successfully."})
	}
}

// ReactivateUser enables a previously soft-deleted user (admin-only workflow)
func (uctlr *UserController) ReactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid user ID format. Please try again."})
	}

	success, err := uctlr.userService.ReactivateUserById(id)
	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Failed to reactivate user, please try again."})
	} else {
		// TODO: propagate an expiry signal to prompt application to trigger password reset dialog
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User reactivated successfully."})
	}
}

// DeleteUser hard-deletes the user from the DB and erases all related data. Used for privacy compliance,
// but not a common workflow.
func (uctlr *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format. Please try again."})
		return
	}
	// TODO: add user exists check and guard

	err = uctlr.userService.DeleteUserById(id)
	if err != nil {
		log.Println("Error deleting user:", err) // TODO add consistent logging with appropriate context throughout
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "Failed to delete user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}
