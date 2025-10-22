package routes	

import(
	"example/vigovia-itenary-api/controllers"
	"example/vigovia-itenary-api/service"
	"example/vigovia-itenary-api/repository"
	"example/vigovia-itenary-api/config"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine , cfg *config.Config){
	//initializes the in memory rep
	repo:=repository.NewInMemoryRepo()
	
	//initializes and creates the itinerary service with repository
	itiSvc:=service.NewItineraryService(repo)

	//initializes and creates the itinerary service using outputDir from Config
	pdfService:=service.NewPDFService(cfg.OutputDir)

	//initializes the route controller with itinerary service and pdfService
	rc:=controllers.NewRouteController(itiSvc,pdfService)

	//sets up the api version group
	v1:=router.Group("/api/v1")
	{
		//group for itinerary related routes
		itineraries:=v1.Group("/itineraries")
		{
			itineraries.POST("",rc.CreateItinerary) //create a new itinerary
			itineraries.GET("",rc.GetAllItineraries) // get all the itineraries
			itineraries.GET("/:id",rc.GetItinerary)  // get itinerary by id
			itineraries.PUT("/:id",rc.UpdateItinerary) //update itinerary
			itineraries.DELETE("/:id",rc.DeleteItinerary) //delete itinerary by id
			itineraries.POST("/:id/pdf", rc.GeneratePDF) //generate pdf for an itinerary by id
			itineraries.GET("/:id/pdf/download", rc.DownloadPDF)  //downloading the pdf for the itinerary
		}
	}

	router.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status":"healthy",
		})
	})
}