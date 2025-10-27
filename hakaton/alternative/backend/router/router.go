package router

import (
	"net/http"

	"backend/apiutils"
	"backend/clients"
	"backend/config"
	"backend/delivery"
	"backend/middleware"
	"backend/repository"
	"backend/services"
	"backend/store"
	"backend/usecase"

	"github.com/gorilla/mux"
	_ "backend/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	redisStore *store.RedisStore,
	minioStore *store.MinIOStore,
	repo *repository.Repository,
	pythonClient *clients.PythonServiceClient,
	cfg *config.Config,
) *mux.Router {

	emailService := services.NewEmailService(&cfg.SMTP)
	userUsecase := usecase.NewUserUsecase(repo, redisStore, minioStore, emailService)
	cometUsecase := usecase.NewCometUsecase(repo)
	observationUsecase := usecase.NewObservationUsecase(repo, minioStore)
	calculationUsecase := usecase.NewCalculationUsecase(repo, pythonClient)

	userHandler := delivery.NewUserHandler(userUsecase, minioStore)
	cometHandler := delivery.NewCometHandler(cometUsecase)
	observationHandler := delivery.NewObservationHandler(observationUsecase)
	calculationHandler := delivery.NewCalculationHandler(calculationUsecase)

	router := mux.NewRouter()
	// router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.CORS)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		apiutils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods("GET")

	// Redirect reset-password to static page
	router.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/reset-password.html")
	})
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
	authRouter.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	authRouter.HandleFunc("/forgot-password", userHandler.RequestPasswordReset).Methods("POST")
	authRouter.HandleFunc("/reset-password", userHandler.ResetPassword).Methods("POST")
	router.HandleFunc("/api/files/{filepath:.*}", userHandler.GetAvatar).Methods("GET")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(redisStore))

	protected.HandleFunc("/profile/me", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile/me", userHandler.UpdateProfile).Methods("PUT")

	protected.HandleFunc("/comets", cometHandler.CreateComet).Methods("POST")
	protected.HandleFunc("/comets", cometHandler.ListUserComets).Methods("GET")
	protected.HandleFunc("/comets/{id:[0-9]+}", cometHandler.GetComet).Methods("GET")
	protected.HandleFunc("/comets/{id:[0-9]+}", cometHandler.UpdateComet).Methods("PUT")
	protected.HandleFunc("/comets/{id:[0-9]+}", cometHandler.DeleteComet).Methods("DELETE")

	protected.HandleFunc("/comets/{comet_id:[0-9]+}/observations", observationHandler.CreateObservationForComet).Methods("POST")
	protected.HandleFunc("/comets/{comet_id:[0-9]+}/observations", observationHandler.ListObservationsForComet).Methods("GET")
	protected.HandleFunc("/observations/{observation_id:[0-9]+}", observationHandler.GetObservation).Methods("GET")
	protected.HandleFunc("/observations/{observation_id:[0-9]+}", observationHandler.UpdateObservation).Methods("PUT")
	protected.HandleFunc("/observations/{observation_id:[0-9]+}", observationHandler.DeleteObservation).Methods("DELETE")

	protected.HandleFunc("/comets/{comet_id:[0-9]+}/calculations", calculationHandler.CreateOrbitCalculation).Methods("POST")
	protected.HandleFunc("/calculations/{calculation_id:[0-9]+}/approach", calculationHandler.AddCloseApproachData).Methods("POST")
	
	protected.HandleFunc("/comets/{comet_id:[0-9]+}/calculations", calculationHandler.ListCalculationsForComet).Methods("GET")
	protected.HandleFunc("/calculations/{calculation_id:[0-9]+}", calculationHandler.GetCalculation).Methods("GET")
	protected.HandleFunc("/calculations/{calculation_id:[0-9]+}", calculationHandler.DeleteCalculation).Methods("DELETE")

	return router
}
