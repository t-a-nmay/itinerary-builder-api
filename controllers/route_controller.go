package controllers

import (
	"example/vigovia-itenary-api/models"
	"example/vigovia-itenary-api/service"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

//acts as a handler for HTTP Requests related to itineraries
type RouteController struct {
	service *service.ItineraryService
	pdfService *service.PDFService
}

//NewRouteController acts as a constructor for RouteController and creates and returns a new RouteController Instance
func NewRouteController(s *service.ItineraryService, pdfSvc *service.PDFService ) *RouteController {
	return &RouteController{
		service: s,
		pdfService: pdfSvc,
	}
}

//CreateItinerary handles POST /api/itineraries
//binds incoming json to the CreateItinerary service and delegates the task to CreateItinerary service
func (rc *RouteController) CreateItinerary(c *gin.Context){
	var req models.CreateItineraryReq
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itinerary,err:=rc.service.CreateItinerary(&req)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, itinerary)
}

// GetItinerary handles GET /api/itineraries/:id
//gets itineraries by id
func (rc *RouteController) GetItinerary(c *gin.Context) {
	id := c.Param("id")

	//asks the GetItinerary from service to fetch the itinerary
	itinerary, err := rc.service.GetItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Itinerary not found",
		})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// GetAllItineraries handles GET /api/itineraries
//gets all the itineraries
func (h *RouteController) GetAllItineraries(c *gin.Context) {
	//filters by user_id
	userID := c.Query("user_id")

	var itineraries []*models.Itinerary
	var err error

	if userID != "" {
		//calls GetUserItineraries from service to get itineraries filtered by user_id
		itineraries, err = h.service.GetUserItineraries(userID)
	} else {
		//calls GetAllItineraries from service to get all the itineraries in case user_id is null
		itineraries, err = h.service.GetAllItineraries()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve itineraries",
		})
		return
	}

	c.JSON(http.StatusOK, itineraries)
}

// UpdateItinerary handles PUT /api/itineraries/:id
func (rc *RouteController) UpdateItinerary(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateItineraryReq

	//binds json to update struct UpdateItineraryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	//calls UpdateItinerary from service to update the itinerary
	itinerary, err := rc.service.UpdateItinerary(id, &req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "itinerary not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// DeleteItinerary handles DELETE /api/itineraries/:id
//deletes the itinerary
func (rc *RouteController) DeleteItinerary(c *gin.Context) {
	id := c.Param("id")

	//calls DeleteItinerary from service
	if err := rc.service.DeleteItinerary(id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "itinerary not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Itinerary deleted successfully",
	})
}


func (rc *RouteController) GeneratePDF(c *gin.Context) {
	id := c.Param("id")

	//calls GetItinerary to fetch the itinerary
	itinerary, err := rc.service.GetItinerary(id)
	if itinerary==nil {
		fmt.Println("cant find itinerary")
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":err.Error(),
		})
		return
	}

	// calls GeneratePDF from pdfService to generate pdf
	filepath, err := rc.pdfService.GeneratePDF(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
	"message": "PDF generated successfully",
	"data": map[string]interface{}{
		"filepath": filepath,
		"filename": filepath[len(filepath)-48:], // Or use path.Base(filepath)
	},
})
}

// DownloadPDF handles GET /api/itineraries/:id/pdf/download
func (rc *RouteController) DownloadPDF(c *gin.Context) {
	id := c.Param("id")

	//calls GetItinerary to fetch the itinerary
	itinerary, err := rc.service.GetItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":err.Error(),
		})
		return
	}

	// calls GeneratePDF from pdfService to generate pdf
	filepath, err := rc.pdfService.GeneratePDF(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}

	// Serve PDF file 
	c.FileAttachment(filepath, filepath[len(filepath)-48:])
}
