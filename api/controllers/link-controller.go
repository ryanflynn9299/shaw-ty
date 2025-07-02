package controllers

import (
	"URLShortener/api/dto"
	id_generator "URLShortener/internal/core/id-generator"
	"URLShortener/internal/services"
	"URLShortener/internal/storage/models"
	"URLShortener/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LinkController struct {
	linkService *services.LinkService
	idGenerator *id_generator.SnowflakeGenerator
}

// TODO: Finish implementing this file

// NewLinkController initializes a new LinkController
func NewLinkController(linkService *services.LinkService) LinkController {
	return LinkController{
		linkService: linkService,
	}
}

// CreateLink creates a new shortlink entry
//
//	> services POST /short_link
func (lctr *LinkController) CreateLink(c *gin.Context) {
	var request dto.CreateLinkRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	shortName, err := lctr.linkService.CreateLink(request.URL, 0, *request.ShortCode, 0, lctr.idGenerator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusCreated, gin.H{"link": shortName, "msg": "Successfully created link with short code " + shortName})
	}
}

// GetLink retrieves a given link by its ID
func (lctr *LinkController) GetLink(c *gin.Context) {
	id := c.Param("id")
	link, err := lctr.linkService.GetLinkById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, shortLinkFromModelCvtr(link))
}

// GetFullLink retireves a shortlink's full URL for redirection
func (lctr *LinkController) GetFullLink(c *gin.Context) {
	id := c.Param("code")
	link, err := lctr.linkService.GetLinkByCode(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"link": link.FullURL})
	}
}

// UpdateLink modifies the properties of a shortlink
func (lctr *LinkController) UpdateLink(c *gin.Context) {

}

// DeactivateLink soft-deletes the short-link by marking it inactive
func (lctr *LinkController) DeactivateLink(c *gin.Context) {
	id := c.Param("id")
	err := lctr.linkService.DeactivateLink(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "Link successfully deactivated."})
	}

}

// DeleteLink hard-deletes the shortlink, permanently removing it from the database
func (lctr *LinkController) DeleteLink(c *gin.Context) {
	id := c.Param("id")
	err := lctr.linkService.DeleteLink(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "Link successfully deleted."})
	}
}

// Helpers:
func shortLinkFromModelCvtr(model *models.ShortLink) dto.LinkResponse {
	return dto.LinkResponse{
		ID:        model.ID,
		ShortCode: model.ShortenedCode,
		URL:       model.FullURL,
		CreatedAt: utils.ConvertUnixToLocalDateString(model.CreatedDate),
		Active:    model.IsActive,
	}
}
