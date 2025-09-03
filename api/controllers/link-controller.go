package controllers

import (
	"URLShortener/api/dto"
	"URLShortener/internal/i18n"
	"URLShortener/internal/services"
	"URLShortener/internal/storage/models"
	"URLShortener/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TODO: add input sanitization and validation

type LinkController struct {
	linkService *services.LinkService
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

	// TODO: generify this error msg
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// retrieve userId from session; another endpoint allows for creation of links without user id
	userId, exists := c.Get("userID")
	if exists == false || userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login and try again."})
		return
	}

	// TODO add userId and expires after
	shortName, err := lctr.linkService.CreateLink(request.URL, int64(userId.(uint)), request.ShortCode, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusCreated, gin.H{"link": utils.GetBaseURL() + shortName, "msg": "Successfully created link with short code " + shortName})
	}
}

// GetLink retrieves a given link by its ID
func (lctr *LinkController) GetLink(c *gin.Context) {
	id := c.Param("id")
	link, err := lctr.linkService.GetLinkById(id)

	// TODO: generify this error msg
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, shortLinkFromModelCvtr(link))
}

// GetFullLink retireves a shortlink's full URL for redirection
func (lctr *LinkController) GetFullLink(c *gin.Context) {
	id := c.Param("code") // TODO: what happens when Param is null?
	shouldRedirect := c.Param("redirect")

	link, err := lctr.linkService.GetLinkByCode(id)
	if err != nil {
		// an error occurred, throw a generic bad JSON response
		respondBadJSONRequest(c)
	} else if shouldRedirect == "false" {
		// a valid link request successfully returned a URL
		c.JSON(http.StatusOK, gin.H{"link": link.FullURL, "status": "success"})
	} else {
		// Redirect request for full URL
		c.Redirect(http.StatusFound, link.FullURL)
	}
}

// GetAllLinksByUser retrieves a given link by its ID
func (lctr *LinkController) GetAllLinksByUser(c *gin.Context) {
	idStr := c.Query("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid user Id. Make sure the user_id parameter is correct and try again."})
		return
	}
	// TODO: add permission check that asks if logged in user has access to the other

	links, err := lctr.linkService.GetAllLinksByUser(int64(id))

	// TODO: generify this error msg
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, links)
}

// UpdateLink modifies the properties of a shortlink
func (lctr *LinkController) UpdateLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID format"})
		return
	}

	// TODO: add msgs to string yaml file
	// TODO: defer dto conversion to service
	var request dto.UpdateLinkRequest
	err = c.ShouldBindJSON(&request) // parse request
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Could not parse JSON body. Please check your request body and try again."})
	}

	// assess RBAP TODO
	// load requested data from db and validate request id
	linkData, err := lctr.linkService.GetLinkById(c.Param("id"))
	if err != nil || linkData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "The link id provided is invalid, check the link id in your request and try again."})
	}

	// make changes
	success, err := lctr.linkService.UpdateLink(id, *request.ShortCode, *request.URL, *request.Active)
	if err != nil || success == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	// post response
	c.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Shortlink updated successfully."})
}

// DeactivateLink soft-deletes the short-link by marking it inactive
func (lctr *LinkController) DeactivateLink(c *gin.Context) {
	id := c.Param("id")
	err := lctr.linkService.DeactivateLink(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Link successfully deactivated."})
	}

}

// DeleteLink hard-deletes the shortlink, permanently removing it from the database
func (lctr *LinkController) DeleteLink(c *gin.Context) {
	id := c.Param("id")
	err := lctr.linkService.DeleteLink(id)
	if err != nil {
		// TODO: make sure not to propagate error message to users OWASP
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Link successfully deleted."})
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

// error response funcs
func respondBadJSONRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest,
		gin.H{"error": i18n.T(i18n.FromAcceptLanguage(c.GetHeader("Accept-Language")), "errors.json_parse_error")})
}
