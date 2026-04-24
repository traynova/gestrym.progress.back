package routes

import (
	"gestrym-progress/docs"
	"gestrym-progress/src/common/middleware"
	"gestrym-progress/src/common/utils"
	progressAdapters "gestrym-progress/src/progress/infrastructure/adapters"
	progressRepos "gestrym-progress/src/progress/infrastructure/repositories"
	progressHandlers "gestrym-progress/src/progress/interfaces/http/handlers"
	progressUseCases "gestrym-progress/src/progress/application/usecases"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gestrym-progress/src/common/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
	metricsHandler *progressHandlers.MetricsHandler
	photosHandler  *progressHandlers.PhotosHandler
	notesHandler   *progressHandlers.NotesHandler
	comparisonHandler *progressHandlers.ComparisonHandler
	workoutHandler    *progressHandlers.WorkoutProgressHandler
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.Title = "Gestrym Progress API"
		docs.SwaggerInfo.Description = "API para el manejo de progresos."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/gestrym-progress"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	// Aplica el middleware CORS
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	// Add default routes
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB
	dbConn := config.NewPostgresConnection()
	db := dbConn.GetDB()

	// Repositories
	metricsRepo := progressRepos.NewGORMBodyMetricsRepository(db)
	photosRepo := progressRepos.NewGORMProgressPhotoRepository(db)
	notesRepo := progressRepos.NewGORMCoachNoteRepository(db)
	workoutRepo := progressRepos.NewGORMWorkoutProgressRepository(db)

	// Adapters & Services
	storageAdapter := progressAdapters.NewStorageServiceAdapter()

	// Use Cases
	createMetricsUC := progressUseCases.NewCreateBodyMetricsUseCase(metricsRepo)
	getMetricsUC := progressUseCases.NewGetUserMetricsUseCase(metricsRepo)
	getWeightChartUC := progressUseCases.NewGetWeightChartUseCase(metricsRepo)
	uploadPhotoUC := progressUseCases.NewUploadProgressPhotoUseCase(photosRepo, storageAdapter)
	getPhotosUC := progressUseCases.NewGetUserPhotosUseCase(photosRepo)
	createNoteUC := progressUseCases.NewCreateCoachNoteUseCase(notesRepo)
	getNotesUC := progressUseCases.NewGetUserNotesUseCase(notesRepo)
	getComparisonUC := progressUseCases.NewGetProgressComparisonUseCase(metricsRepo, photosRepo)
	markWorkoutUC := progressUseCases.NewMarkWorkoutProgressUseCase(workoutRepo)
	getWorkoutUC := progressUseCases.NewGetWorkoutProgressUseCase(workoutRepo)

	// Controllers
	r.metricsHandler = progressHandlers.NewMetricsHandler(createMetricsUC, getMetricsUC, getWeightChartUC)
	r.photosHandler = progressHandlers.NewPhotosHandler(uploadPhotoUC, getPhotosUC)
	r.notesHandler = progressHandlers.NewNotesHandler(createNoteUC, getNotesUC)
	r.comparisonHandler = progressHandlers.NewComparisonHandler(getComparisonUC)
	r.workoutHandler = progressHandlers.NewWorkoutProgressHandler(markWorkoutUC, getWorkoutUC)

	// Add server group
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add groups
	r.publicGroup = r.serverGroup.Group("/public")
	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Add middleware to private group
	r.privateGroup.Use(middleware.SetupJWTMiddleware())

	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	// Add routes to remaining groups
	r.addPrivateRoutes()
	r.addPublicRoutes()
	r.addInternalRoutes()
	r.addProtectedRoutes()

}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {

	// Handle root
	serverInstance.GET("/", func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-progress OK...",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusOK, response)
	})

	// Handle 404
	serverInstance.NoRoute(func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusNotFound, response)
	})
}

func (r *routesDefinition) addPublicRoutes() {

}

func (r *routesDefinition) addPrivateRoutes() {
	// Metrics routes
	metrics := r.privateGroup.Group("/metrics")
	{
		metrics.POST("", r.metricsHandler.Create)
		metrics.GET("/user/:id", r.metricsHandler.GetByUserID)
		metrics.GET("/user/:id/chart", r.metricsHandler.GetWeightChart)
	}

	// Photos routes
	photos := r.privateGroup.Group("/photos")
	{
		photos.POST("", r.photosHandler.Upload)
		photos.GET("/user/:id", r.photosHandler.GetByUserID)
	}

	// Notes routes
	notes := r.privateGroup.Group("/notes")
	{
		notes.POST("", middleware.RequireRoles(middleware.RoleCoach, middleware.RoleGym, middleware.RoleAdmin), r.notesHandler.Create)
		notes.GET("/user/:id", r.notesHandler.GetByUserID)
	}

	// Comparison route
	r.privateGroup.GET("/comparison/user/:id", r.comparisonHandler.GetComparison)

	// Workout Progress routes
	workout := r.privateGroup.Group("/workout-progress")
	{
		workout.POST("", r.workoutHandler.Mark)
		workout.GET("/user/:id", r.workoutHandler.GetByUserID)
	}
}

func (r *routesDefinition) addInternalRoutes() {

}

func (r *routesDefinition) addProtectedRoutes() {
}
