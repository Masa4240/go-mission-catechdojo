package router

import (
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/controller"
	"github.com/Masa4240/go-mission-catechdojo/handler"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/Masa4240/go-mission-catechdojo/service"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

func NewRouter(userDB *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)
	// svc := service.NewUserService(userDB)
	// userHandler := handler.NewUserHandler(svc)

	// r.Route("/user", func(r chi.Router) {
	// 	r.Post("/create", userHandler.CreateUser)
	// 	r.Route("/", func(r chi.Router) {
	// 		r.Use(middleware.TokenValidation)
	// 		r.Get("/get", userHandler.GetUserName)
	// 		r.Put("/update", userHandler.UpdateUser)
	// 	})
	// })

	svc3 := model.NewUserModel(userDB)
	svc2 := service.NewUserServiceMVC(svc3)
	userController := controller.NewUserController(svc2)
	r.Route("/user1", func(r chi.Router) {
		r.Post("/create", userController.CreateUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.TokenValidation)
			r.Get("/get", userController.GetUser)
			r.Put("/update", userController.UpdateUser)
		})
	})

	gsvc3 := model.NewGachaModel(userDB)
	gsvc2 := service.NewGachaServiceSVC(gsvc3)
	gachaController := controller.NewGachaController(gsvc2)
	r.Route("/gacha1", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Post("/draw", gachaController.Gacha)
	})

	r.Route("/character1", func(r chi.Router) {
		r.Use(middleware.TokenValidation)
		r.Get("/list", gachaController.GetCharacterList)
	})

	healthzHandler := handler.NewHealthzHandler()
	r.Get("/healthz", healthzHandler.ServeHTTP)

	// gsvc := service.NewGachaService(userDB)
	// gachaHandler := handler.NewGachaHandler(gsvc)
	// //r.Get("/gacha", gachaHandler.Gacha)
	// r.Route("/gacha", func(r chi.Router) {
	// 	r.Use(middleware.TokenValidation)
	// 	r.Post("/draw", gachaHandler.Gacha)
	// })

	// r.Route("/character", func(r chi.Router) {
	// 	r.Use(middleware.TokenValidation)
	// 	r.Get("/list", gachaHandler.GetCharsList)
	// })

	// For admin Usage
	r.Post("/char/add", gachaController.AddCharacter)

	return r
}
